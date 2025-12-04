package attendance

import (
	"time"

	"github.com/nanato-okajima/attendance_management/internal/entity"
)

//go:generate mockgen -source=repository.go -destination=../../mock/mock_attendance_repository.go -package=mock -mock_names Repository=MockAttendanceRepository

// Repository はAttendanceエンティティの永続化を抽象化
type Repository interface {
	Create(attendance *entity.Attendance) error
	FindByEmployeeAndDate(employeeID int, date time.Time) (*entity.Attendance, error)
	Update(attendance *entity.Attendance) error
	FindByEmployeeAndDateRange(employeeID int, startDate, endDate time.Time) ([]entity.Attendance, error)
}
