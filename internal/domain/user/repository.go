package user

import (
	"github.com/nanato-okajima/attendance_management/internal/entity"
)

//go:generate mockgen -source=repository.go -destination=../../mock/mock_user_repository.go -package=mock -mock_names Repository=MockUserRepository

// Repository はUserエンティティの永続化を抽象化
type Repository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByEmployeeNumber(employeeNumber int) (*entity.User, error)
	UpdateLastLogin(id uint) error
}
