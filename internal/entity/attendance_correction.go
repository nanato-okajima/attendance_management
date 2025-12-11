package entity

import "time"

// CorrectionType 修正種別
type CorrectionType int

const (
	CorrectionTypeClockIn  CorrectionType = 1 // 出勤修正
	CorrectionTypeClockOut CorrectionType = 2 // 退勤修正
	CorrectionTypeBoth     CorrectionType = 3 // 両方修正
)

// AttendanceCorrection 打刻修正申請エンティティ
type AttendanceCorrection struct {
	ID                   uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	AttendanceID         uint           `gorm:"not null;index" json:"attendance_id"`
	EmployeeNumber       int            `gorm:"not null;index" json:"employee_number"`
	CorrectionType       CorrectionType `gorm:"not null" json:"correction_type"`
	CorrectedOpeningTime *time.Time     `json:"corrected_opening_time,omitempty"`
	CorrectedClosingTime *time.Time     `json:"corrected_closing_time,omitempty"`
	Reason               string         `gorm:"type:text;not null" json:"reason"`
	ApprovalStatus       ApprovalStatus `gorm:"default:1;index" json:"approval_status"`
	ApproverID           *int           `json:"approver_id,omitempty"`
	ApprovedAt           *time.Time     `json:"approved_at,omitempty"`
	RejectReason         string         `gorm:"type:text" json:"reject_reason,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (AttendanceCorrection) TableName() string {
	return "attendance_corrections"
}
