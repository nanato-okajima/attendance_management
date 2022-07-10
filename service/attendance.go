package service

import (
	"github.com/nanato-okajima/attendance_management/database"
)

type Attendance struct {
	AttendanceId     int64  `json:"attendance_id"`
	EmployeeId       int64  `json:"employee_id"`
	OpeningTime      string `json:"opening_time"`
	ClosingTime      string `json:"closing_time"`
	AttendanceStatus int64  `json:"attendance_status"`
}

func Register(attendance *Attendance) error {
	_ = database.DB.Client.Create(attendance)

	return nil
}

func List() (*[]Attendance, error) {
	attendances := []Attendance{}
	_ = database.DB.Client.Find(&attendances)

	return &attendances, nil
}

func Update(attendance *Attendance, id string) error {
	database.DB.Client.Model(&Attendance{}).Where("attendance_id = " + id).Updates(attendance)

	return nil
}

func Delete(id string) error {
	database.DB.Client.Where("attendance_id = " + id).Delete(&Attendance{})

	return nil
}
