package entity

import "time"

type Department struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code                  string    `gorm:"uniqueIndex;not null" json:"code"`
	Name                  string    `gorm:"not null" json:"name"`
	ParentID              *uint     `json:"parent_id"`
	ManagerEmployeeNumber *int      `json:"manager_employee_number"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
