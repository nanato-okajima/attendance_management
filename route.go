package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/nanato-okajima/attendance_management/validator"
)

const (
	dsnFormat = "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
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

type DBEnv struct {
	User     string
	Password string
	Host     string
	Name     string
}

func settingRoute() {
	http.HandleFunc("/attendance/", AttendanceHandler)
	http.HandleFunc("/attendance/register", AttendanceRegisterHandler)
	http.HandleFunc("/attendance/list", AttendanceListHandler)
}

func AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/attendance/")

	errMsg, err := validator.Validation(req)
	if err != nil {
		log.Println("error")
		return
	}
	if len(errMsg) > 0 {
		fmt.Println(strings.Join(errMsg, ":"))
		return
	}

	db, err := setupDB()
	if err != nil {
		log.Println("error")
		return
	}

	if r.Method == "PUT" {
		attendance := createRecord(&req)
		db.Model(&Attendance{}).Where("attendance_id = " + id).Updates(attendance)

		err := createResponse(w, http.StatusOK, "updated!")
		if err != nil {
			log.Println("error")
			return
		}
	}

	if r.Method == "DELETE" {
		db.Where("attendance_id = " + id).Delete(&Attendance{})

		err := createResponse(w, http.StatusOK, "deleted!")
		if err != nil {
			log.Println("error")
			return
		}
	}
}

func AttendanceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	errMsg, err := validator.Validation(req)
	if err != nil {
		log.Println("error")
		return
	}
	if len(errMsg) > 0 {
		fmt.Println(strings.Join(errMsg, ":"))
		return
	}

	db, err := setupDB()
	if err != nil {
		log.Println("error")
		return
	}
	attendance := createRecord(&req)
	_ = db.Create(attendance)
	if err := createResponse(w, http.StatusCreated, "created!"); err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
}

func AttendanceListHandler(w http.ResponseWriter, r *http.Request) {
	db, err := setupDB()
	if err != nil {
		log.Println("error")
		return
	}

	attendances := []Attendance{}
	_ = db.Find(&attendances)

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

func createRecord(req *Request) *Attendance {
	return &Attendance{
		EmployeeId:       1,
		OpeningTime:      req.OpeningTime,
		ClosingTime:      req.ClosingTime,
		AttendanceStatus: 1,
	}
}

func setupDB() (*gorm.DB, error) {
	var env DBEnv
	err := envconfig.Process("db", &env)
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(dsnFormat, env.User, env.Password, env.Host, env.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
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
