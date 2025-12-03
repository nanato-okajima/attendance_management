package service

//go:generate mockgen -source=attendance.go -destination=mock/mock_attendance_repository.go -package=mock

import (
	"time"

	"github.com/nanato-okajima/attendance_management/internal/models"
)

type AttendanceRepository interface {
	Create(attendance *models.Attendance) error
	FindByEmployeeAndDate(employeeID int, date time.Time) (*models.Attendance, error)
	Update(attendance *models.Attendance) error
	FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]models.Attendance, error)
}
