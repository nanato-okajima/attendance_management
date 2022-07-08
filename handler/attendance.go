package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nanato-okajima/attendance_management/database"
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

func AttendanceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	if err := validator.Validation(req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	attendance := createRecord(&req)
	_ = database.DB.Client.Create(attendance)
	if err := createResponse(w, http.StatusCreated, "created!"); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
}

func AttendanceListHandler(w http.ResponseWriter, r *http.Request) {
	attendances := []Attendance{}
	_ = database.DB.Client.Find(&attendances)

	e, err := json.Marshal(attendances)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(e)
	if err != nil {
		log.Println("error")
		return
	}
}

func AttendanceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/attendance/")

	if err := validator.Validation(req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	attendance := createRecord(&req)
	database.DB.Client.Model(&Attendance{}).Where("attendance_id = " + id).Updates(attendance)

	if err := createResponse(w, http.StatusOK, "updated!"); err != nil {
		log.Println("error")
		return
	}
}

func AttendanceDeleteHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/attendance/")

	if err := validator.Validation(req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	database.DB.Client.Where("attendance_id = " + id).Delete(&Attendance{})

	if err := createResponse(w, http.StatusOK, "deleted!"); err != nil {
		log.Println("error")
		return
	}
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
