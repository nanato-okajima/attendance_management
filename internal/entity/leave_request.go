package entity

import "time"

type LeaveRequest struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeNumber  int        `gorm:"not null" json:"employee_number"`
	LeaveType       int        `gorm:"not null" json:"leave_type"`
	StartDate       time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate         time.Time  `gorm:"type:date;not null" json:"end_date"`
	HalfDayType     *int       `json:"half_day_type"`
	Reason          string     `gorm:"type:text;not null" json:"reason"`
	ApprovalStatus  int        `gorm:"default:1" json:"approval_status"`
	ApproverID      *int       `json:"approver_id"`
	ApprovedAt      *time.Time `json:"approved_at"`
	ApprovalComment string     `gorm:"type:text" json:"approval_comment"`
	RejectReason    string     `gorm:"type:text" json:"reject_reason"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
