package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	l "github.com/nanato-okajima/attendance_management/logger"
	"github.com/nanato-okajima/attendance_management/myerrors"
	"github.com/nanato-okajima/attendance_management/service"
	"github.com/nanato-okajima/attendance_management/validator"
)

type Request struct {
	OpeningTime string `json:"opening_time" validate:"required,date_format_check"`
	ClosingTime string `json:"closing_time" validate:"required,date_format_check"`
}

type AttendanceHandler struct{}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}

	l.Logger.Infof("attendance register : requests = {opening_time: %s, closing_time: %s}", req.OpeningTime, req.ClosingTime)

	if err := validator.Validation(req); err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}

	attendance := createRecord(&req)

	if err := service.Register(attendance); err != nil {
		errorHandler(w, err)
	}

	createResponse(w, http.StatusCreated, "registered!")
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	attendances, err := service.List()
	if err != nil {
		errorHandler(w, err)
	}

	e, err := json.Marshal(attendances)
	if err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(e)
	if err != nil {
		errorHandler(w, &myerrors.InternalServerError{Err: err})
		return
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}
	id := mux.Vars(r)["id"]

	l.Logger.Infof("attendance register : {id: %s} : requests = {opening_time: %s, closing_time: %s}", id, req.OpeningTime, req.ClosingTime)

	if err := validator.Validation(req); err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}

	attendance := createRecord(&req)

	if err := service.Update(attendance, id); err != nil {
		errorHandler(w, err)
	}
	createResponse(w, http.StatusOK, "updated!")
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}
	id := mux.Vars(r)["id"]

	l.Logger.Infof("attendance register : {id: %s} : requests = {opening_time: %s, closing_time: %s}", id, req.OpeningTime, req.ClosingTime)

	if err := validator.Validation(req); err != nil {
		errorHandler(w, &myerrors.BadRequestError{Err: err})
		return
	}

	if err := service.Delete(id); err != nil {
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

func createRecord(req *Request) *service.Attendance {
	return &service.Attendance{
		EmployeeId:       1,
		OpeningTime:      req.OpeningTime,
		ClosingTime:      req.ClosingTime,
		AttendanceStatus: 1,
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
