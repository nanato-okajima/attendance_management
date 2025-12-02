package database

import (
	"time"

	"github.com/nanato-okajima/attendance_management/models"
	"github.com/nanato-okajima/attendance_management/service/repository"
)

type userRepository struct{}

func NewUserRepository() repository.UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmployeeNumber(employeeNumber int) (*models.User, error) {
	var user models.User
	err := DB.Where("employee_number = ? AND is_active = ?", employeeNumber, true).First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateLastLogin(id uint) error {
	now := time.Now()
	return DB.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}
