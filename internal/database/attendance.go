package database

import (
	"time"

	"github.com/nanato-okajima/attendance_management/internal/domain/attendance"
	"github.com/nanato-okajima/attendance_management/internal/entity"
)

type attendanceRepository struct{}

func NewAttendanceRepository() attendance.Repository {
	return &attendanceRepository{}
}

func (r *attendanceRepository) Create(attendance *entity.Attendance) error {
	return DB.Create(attendance).Error
}

func (r *attendanceRepository) FindByEmployeeAndDate(employeeID int, date time.Time) (*entity.Attendance, error) {
	var attendance entity.Attendance
	err := DB.Where("employee_id = ? AND target_date = ?", employeeID, date.Format("2006-01-02")).First(&attendance).Error
	return &attendance, err
}

func (r *attendanceRepository) Update(attendance *entity.Attendance) error {
	return DB.Save(attendance).Error
}

func (r *attendanceRepository) FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]entity.Attendance, error) {
	var attendances []entity.Attendance
	err := DB.Where("employee_id = ? AND target_date BETWEEN ? AND ?",
		employeeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("target_date ASC").
		Find(&attendances).Error
	return attendances, err
}
