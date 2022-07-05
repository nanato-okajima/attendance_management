package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	http.HandleFunc("/attendance/register", AttendanceRegisterHandler)
}

func AttendanceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	var env DBEnv
	envconfig.Process("db", &env)

	dsn := fmt.Sprintf(dsnFormat, env.User, env.Password, env.Host, env.Name)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err)
	}

	attendance := createRecord(&req)

	_ = db.Create(attendance)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	m := map[string]string{"massage": "created!"}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Println(err)
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
