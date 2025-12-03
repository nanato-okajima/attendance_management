package repository

import "github.com/nanato-okajima/attendance_management/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByEmployeeNumber(employeeNumber int) (*models.User, error)
	UpdateLastLogin(id uint) error
}
