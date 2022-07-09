package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nanato-okajima/attendance_management/database"
	l "github.com/nanato-okajima/attendance_management/logger"
	"github.com/nanato-okajima/attendance_management/myerrors"
	"github.com/nanato-okajima/attendance_management/validator"
)

type Request struct {
	OpeningTime string `json:"opening_time" validate:"required,date_format_check"`
	ClosingTime string `json:"closing_time" validate:"required,date_format_check"`
}

type Attendance struct {
	AttendanceId     int64  `json:"attendance_id"`
	EmployeeId       int64  `json:"employee_id"`
	OpeningTime      string `json:"opening_time"`
	ClosingTime      string `json:"closing_time"`
	AttendanceStatus int64  `json:"attendance_status"`
}

type AttendanceHandler struct{}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := Register(w, r); err != nil {
		errorHandler(w, err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) error {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myerrors.BadRequestError{Err: err}
	}

	l.Logger.Infof("attendance register : requests = {opening_time: %s, closing_time: %s}", req.OpeningTime, req.ClosingTime)

	if err := validator.Validation(req); err != nil {
		return &myerrors.BadRequestError{Err: err}
	}

	attendance := createRecord(&req)
	_ = database.DB.Client.Create(attendance)
	if err := createResponse(w, http.StatusCreated, "created!"); err != nil {
		return &myerrors.InternalServerError{Err: err}
	}

	return nil
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if err := List(w, r); err != nil {
		errorHandler(w, err)
	}
}

func List(w http.ResponseWriter, r *http.Request) error {
	attendances := []Attendance{}
	_ = database.DB.Client.Find(&attendances)

	e, err := json.Marshal(attendances)
	if err != nil {
		return &myerrors.BadRequestError{Err: err}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(e)
	if err != nil {
		return &myerrors.InternalServerError{Err: err}
	}

	return nil
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if err := Update(w, r); err != nil {
		errorHandler(w, err)
	}
}

func Update(w http.ResponseWriter, r *http.Request) error {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myerrors.BadRequestError{Err: err}
	}
	id := mux.Vars(r)["id"]

	l.Logger.Infof("attendance register : {id: %s} : requests = {opening_time: %s, closing_time: %s}", id, req.OpeningTime, req.ClosingTime)

	if err := validator.Validation(req); err != nil {
		return &myerrors.BadRequestError{Err: err}
	}

	attendance := createRecord(&req)
	database.DB.Client.Model(&Attendance{}).Where("attendance_id = " + id).Updates(attendance)

	if err := createResponse(w, http.StatusOK, "updated!"); err != nil {
		return &myerrors.InternalServerError{Err: err}
	}

	return nil
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if err := Delete(w, r); err != nil {
		errorHandler(w, err)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) error {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &myerrors.BadRequestError{Err: err}
	}
	id := mux.Vars(r)["id"]

	l.Logger.Infof("attendance register : {id: %s} : requests = {opening_time: %s, closing_time: %s}", id, req.OpeningTime, req.ClosingTime)

	if err := validator.Validation(req); err != nil {
		return &myerrors.BadRequestError{Err: err}
	}

	database.DB.Client.Where("attendance_id = " + id).Delete(&Attendance{})

	if err := createResponse(w, http.StatusOK, "deleted!"); err != nil {
		return &myerrors.InternalServerError{Err: err}
	}

	return nil
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{}
}

func createRecord(req *Request) *Attendance {
	return &Attendance{
		EmployeeId:       1,
		OpeningTime:      req.OpeningTime,
		ClosingTime:      req.ClosingTime,
		AttendanceStatus: 1,
	}
}

func createResponse(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	m := map[string]string{"massage": message}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		return err
	}

	return nil
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
