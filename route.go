package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsnFormat = "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

type Request struct {
	OpeningTime string `json:"opening_time" validate:"required,date_format_check"`
	ClosingTime string `json:"closing_time" validate:"required,date_format_check"`
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

	errMsg := validation(req)
	if len(errMsg) > 0 {
		fmt.Println(strings.Join(errMsg, ":"))
		return
	}

	db := setupDB()

	if r.Method == "PUT" {
		attendance := createRecord(&req)
		db.Model(&Attendance{}).Where("attendance_id = " + id).Updates(attendance)

		createResponse(w, http.StatusOK, "updated!")
	}

	if r.Method == "DELETE" {
		db.Where("attendance_id = " + id).Delete(&Attendance{})

		createResponse(w, http.StatusOK, "deleted!")
	}
}

func AttendanceRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	errMsg := validation(req)
	if len(errMsg) > 0 {
		fmt.Println(strings.Join(errMsg, ":"))
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

func validation(attendance Request) []string {
	validate := validator.New()
	validate.RegisterValidation("date_format_check", DateFormatCheck)

	if err := validate.Struct(attendance); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			var errorMessage string
			fieldName := err.Field()

			switch fieldName {
			case "OpeningTime":
				typ := err.Tag()
				switch typ {
				case "required":
					errorMessage = "OpeningTimeは必須項目です"
				case "date_format_check":
					errorMessage = "OpeningTimeの日付形式が不正です"
				}
			case "ClosingTime":
				typ := err.Tag()
				switch typ {
				case "required":
					errorMessage = "OpeningTimeは必須項目です"
				case "date_format_check":
					errorMessage = "ClosingTimeの日付形式が不正です"
				}
			}
			errorMessages = append(errorMessages, errorMessage)
			return errorMessages
		}
	}
	return nil
}

func DateFormatCheck(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02 15:04:05", fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
