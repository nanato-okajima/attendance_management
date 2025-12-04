package entity

import "time"

type User struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	EmployeeNumber int        `gorm:"not null" json:"employee_number"`
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string     `gorm:"not null" json:"-"`
	Role           string     `gorm:"type:enum('admin','employee');default:'employee'" json:"role"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
