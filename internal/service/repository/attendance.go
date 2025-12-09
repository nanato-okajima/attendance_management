package service

//go:generate mockgen -source=attendance.go -destination=../..mock/mock_attendance_repository.go -package=mock

import (
	"time"

	"github.com/nanato-okajima/attendance_management/internal/entity"
)

type AttendanceRepository interface {
	Create(attendance *entity.Attendance) error
	FindByEmployeeAndDate(employeeID int, date time.Time) (*entity.Attendance, error)
	Update(attendance *entity.Attendance) error
	FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]entity.Attendance, error)
}
