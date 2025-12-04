package entity

import "time"

type LineLinkingCode struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeNumber int        `gorm:"not null" json:"employee_number"`
	Code           string     `gorm:"uniqueIndex;not null" json:"code"`
	ExpiresAt      time.Time  `gorm:"not null;index" json:"expires_at"`
	IsUsed         bool       `gorm:"default:false" json:"is_used"`
	UsedAt         *time.Time `json:"used_at"`
	CreatedAt      time.Time  `json:"created_at"`
}
