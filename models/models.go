package models

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

type Department struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code                  string    `gorm:"uniqueIndex;not null" json:"code"`
	Name                  string    `gorm:"not null" json:"name"`
	ParentID              *uint     `json:"parent_id"`
	ManagerEmployeeNumber *int      `json:"manager_employee_number"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type Position struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"not null" json:"name"`
	Level     int       `gorm:"not null" json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LineUser struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeNumber int       `gorm:"not null;uniqueIndex" json:"employee_number"`
	LineUserID     string    `gorm:"uniqueIndex;not null" json:"line_user_id"`
	LinkedAt       time.Time `gorm:"not null" json:"linked_at"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LineLinkingCode struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeNumber int        `gorm:"not null" json:"employee_number"`
	Code           string     `gorm:"uniqueIndex;not null" json:"code"`
	ExpiresAt      time.Time  `gorm:"not null;index" json:"expires_at"`
	IsUsed         bool       `gorm:"default:false" json:"is_used"`
	UsedAt         *time.Time `json:"used_at"`
	CreatedAt      time.Time  `json:"created_at"`
}
