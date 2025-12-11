package entity

import "time"

// LeaveType 休暇種別
type LeaveType int

const (
	LeaveTypePaidLeave    LeaveType = 1 // 有給休暇
	LeaveTypeSpecialLeave LeaveType = 2 // 特別休暇
	LeaveTypeAbsence      LeaveType = 3 // 欠勤
	LeaveTypeHalfDay      LeaveType = 4 // 半休
)

// HalfDayType 半休区分
type HalfDayType int

const (
	HalfDayTypeMorning   HalfDayType = 1 // 午前半休
	HalfDayTypeAfternoon HalfDayType = 2 // 午後半休
)

// ApprovalStatus 承認状態
type ApprovalStatus int

const (
	ApprovalStatusPending  ApprovalStatus = 1 // 承認待ち
	ApprovalStatusApproved ApprovalStatus = 2 // 承認済み
	ApprovalStatusRejected ApprovalStatus = 3 // 却下
)

// LeaveRequest 休暇申請エンティティ
type LeaveRequest struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeNumber  int            `gorm:"not null;index" json:"employee_number"`
	LeaveType       LeaveType      `gorm:"not null" json:"leave_type"`
	StartDate       time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate         time.Time      `gorm:"type:date;not null" json:"end_date"`
	HalfDayType     *HalfDayType   `json:"half_day_type,omitempty"`
	Reason          string         `gorm:"type:text;not null" json:"reason"`
	ApprovalStatus  ApprovalStatus `gorm:"default:1;index" json:"approval_status"`
	ApproverID      *int           `json:"approver_id,omitempty"`
	ApprovedAt      *time.Time     `json:"approved_at,omitempty"`
	ApprovalComment string         `gorm:"type:text" json:"approval_comment,omitempty"`
	RejectReason    string         `gorm:"type:text" json:"reject_reason,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (LeaveRequest) TableName() string {
	return "leave_requests"
}
