package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsnFormat = "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

type Request struct {
	OpeningTime string `json:"opening_time"`
	ClosingTime string `json:"closing_time"`
}

type Attendance struct {
	EmployeeId       int64
	OpeningTime      string
	ClosingTime      string
	AttendanceStatus int64
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
}

func AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/attendance/")

	db := setupDB()
	attendance := createRecord(&req)
	db.Model(&Attendance{}).Where("attendance_id = " + id).Updates(attendance)

	createResponse(w, http.StatusOK, "updated!")
}

func AttendanceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	db := setupDB()
	attendance := createRecord(&req)
	_ = db.Create(attendance)
	if err := createResponse(w, http.StatusCreated, "created!"); err != nil {
		fmt.Fprintf(w, err.Error())
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

func setupDB() *gorm.DB {
	var env DBEnv
	envconfig.Process("db", &env)

	dsn := fmt.Sprintf(dsnFormat, env.User, env.Password, env.Host, env.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err)
	}

	return db
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
