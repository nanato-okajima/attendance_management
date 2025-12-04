package entity

import "time"

type Position struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"not null" json:"name"`
	Level     int       `gorm:"not null" json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
