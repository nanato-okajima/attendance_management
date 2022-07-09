package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	l "github.com/nanato-okajima/attendance_management/logger"

	"github.com/nanato-okajima/attendance_management/database"
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

	return nil
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

	return nil
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

	return nil
}

func createRecord(req *Request) *Attendance {
	return &Attendance{
		EmployeeId:       1,
		OpeningTime:      req.OpeningTime,
		ClosingTime:      req.ClosingTime,
		AttendanceStatus: 1,
	}
}
