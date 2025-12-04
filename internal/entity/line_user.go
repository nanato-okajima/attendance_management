package entity

import "time"

type LineUser struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeNumber int       `gorm:"not null;uniqueIndex" json:"employee_number"`
	LineUserID     string    `gorm:"uniqueIndex;not null" json:"line_user_id"`
	LinkedAt       time.Time `gorm:"not null" json:"linked_at"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
