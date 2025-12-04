package database

import (
	"time"

	"github.com/nanato-okajima/attendance_management/internal/domain/user"
	"github.com/nanato-okajima/attendance_management/internal/entity"
)

type userRepository struct{}

func NewUserRepository() user.Repository {
	return &userRepository{}
}

func (r *userRepository) Create(user *entity.User) error {
	return DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := DB.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmployeeNumber(employeeNumber int) (*entity.User, error) {
	var user entity.User
	err := DB.Where("employee_number = ? AND is_active = ?", employeeNumber, true).First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	now := time.Now()
	return DB.Model(&entity.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}
