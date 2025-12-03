package database

import (
	"time"

	"github.com/nanato-okajima/attendance_management/internal/domain/service"
	"github.com/nanato-okajima/attendance_management/internal/models"
)

type attendanceRepository struct{}

func NewAttendanceRepository() service.AttendanceRepository {
	return &attendanceRepository{}
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	return DB.Create(attendance).Error
}

func (r *attendanceRepository) FindByEmployeeAndDate(employeeID int, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	err := DB.Where("employee_id = ? AND target_date = ?", employeeID, date.Format("2006-01-02")).First(&attendance).Error
	return &attendance, err
}

func (r *attendanceRepository) Update(attendance *models.Attendance) error {
	return DB.Save(attendance).Error
}

func (r *attendanceRepository) FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := DB.Where("employee_id = ? AND target_date BETWEEN ? AND ?",
		employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("target_date ASC").
		Find(&attendances).Error
	return attendances, err
}
