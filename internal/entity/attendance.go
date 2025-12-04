package entity

import "time"

type Attendance struct {
	AttendanceID     uint       `gorm:"primaryKey;autoIncrement" json:"attendance_id"`
	EmployeeID       int        `gorm:"not null" json:"employee_id"`
	TargetDate       time.Time  `gorm:"type:date;not null" json:"target_date"`
	OpeningTime      *time.Time `json:"opening_time"`
	ClosingTime      *time.Time `json:"closing_time"`
	AttendanceStatus int        `gorm:"not null" json:"attendance_status"`
	ClockSource      int        `gorm:"default:1" json:"clock_source"`
	Latitude         *float64   `json:"latitude"`
	Longitude        *float64   `json:"longitude"`
	WorkHours        *float64   `json:"work_hours"`
	OvertimeHours    *float64   `json:"overtime_hours"`
	Note             string     `json:"note"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
