package database

import (
	"github.com/nanato-okajima/attendance_management/service/repository"
)

type attendanceDatabase struct {
}

func NewAttendanceDatabase() repository.AttendanceRepository {
	return &attendanceDatabase{}
}

func (ad attendanceDatabase) Insert(attendance *repository.Attendance) error {
	result := DB.Client.Create(attendance)

	return result.Error
}

func (ad attendanceDatabase) Fetch(attendances *[]repository.Attendance) error {
	result := DB.Client.Find(&attendances)

	return result.Error
}

func (ad attendanceDatabase) Update(attendance *repository.Attendance, id string) error {
	result := DB.Client.Model(&repository.Attendance{}).Where("attendance_id = " + id).Updates(attendance)

	return result.Error
}

func (ad attendanceDatabase) Delete(id string) error {
	result := DB.Client.Where("attendance_id = " + id).Delete(&repository.Attendance{})
	return result.Error
}
