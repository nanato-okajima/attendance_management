package service

import "github.com/nanato-okajima/attendance_management/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByEmployeeNumber(employeeNumber int) (*entity.User, error)
	UpdateLastLogin(id uint) error
}
