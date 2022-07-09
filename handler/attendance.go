package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	l "github.com/nanato-okajima/attendance_management/logger"
	"github.com/nanato-okajima/attendance_management/myerrors"
	"github.com/nanato-okajima/attendance_management/service"
)

type AttendanceHandler struct{}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := service.Register(w, r); err != nil {
		errorHandler(w, err)
	}

	createResponse(w, http.StatusCreated, "registered!")
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if err := service.List(w, r); err != nil {
		errorHandler(w, err)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := service.Update(w, r); err != nil {
		errorHandler(w, err)
	}

	createResponse(w, http.StatusOK, "updated!")
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if err := service.Delete(w, r); err != nil {
		errorHandler(w, err)
	}

	createResponse(w, http.StatusOK, "deleted!")
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{}
}

func errorHandler(w http.ResponseWriter, err error) {
	var br *myerrors.BadRequestError
	if errors.As(err, &br) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		l.Logger.Errorf("400 ", err)
	}

	var is *myerrors.InternalServerError
	if errors.As(err, &is) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		l.Logger.Errorf("500 ", err)
	}
}

func createResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	m := map[string]string{"massage": message}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
