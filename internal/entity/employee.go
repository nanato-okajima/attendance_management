package entity

import "time"

type Employee struct {
	EmployeeNumber int        `gorm:"primaryKey;autoIncrement" json:"employee_number"`
	Name           string     `gorm:"not null" json:"name"`
	NameKana       string     `gorm:"not null" json:"name_kana"`
	Birthday       time.Time  `gorm:"type:date;not null" json:"birthday"`
	GenderCd       int        `gorm:"not null" json:"gender_cd"`
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	Phone          string     `json:"phone"`
	DepartmentID   *int       `json:"department_id"`
	PositionID     *int       `json:"position_id"`
	HireDate       time.Time  `gorm:"type:date;not null" json:"hire_date"`
	EmploymentType int        `gorm:"not null" json:"employment_type"`
	IsDeleted      bool       `gorm:"default:false" json:"is_deleted"`
	DeletedAt      *time.Time `json:"deleted_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
