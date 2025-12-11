package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
	"github.com/nanato-okajima/attendance_management/internal/service"
)

// LeaveHandler 休暇申請ハンドラー
type LeaveHandler struct {
	leaveService service.LeaveService
}

// NewLeaveHandler 休暇申請ハンドラーを作成
func NewLeaveHandler(leaveService service.LeaveService) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leaveService,
	}
}

// CreateLeaveRequestRequest 休暇申請作成リクエスト
type CreateLeaveRequestRequest struct {
	LeaveType   int    `json:"leave_type" binding:"required,min=1,max=4"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	HalfDayType *int   `json:"half_day_type,omitempty"`
	Reason      string `json:"reason" binding:"required,max=500"`
}

// LeaveRequestResponse 休暇申請レスポンス
type LeaveRequestResponse struct {
	ID              uint      `json:"id"`
	EmployeeNumber  int       `json:"employee_number"`
	LeaveType       int       `json:"leave_type"`
	StartDate       string    `json:"start_date"`
	EndDate         string    `json:"end_date"`
	HalfDayType     *int      `json:"half_day_type,omitempty"`
	Reason          string    `json:"reason"`
	ApprovalStatus  int       `json:"approval_status"`
	ApproverID      *int      `json:"approver_id,omitempty"`
	ApprovedAt      *string   `json:"approved_at,omitempty"`
	ApprovalComment string    `json:"approval_comment,omitempty"`
	RejectReason    string    `json:"reject_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateLeaveRequest 休暇申請作成
func (h *LeaveHandler) CreateLeaveRequest(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	var req CreateLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("開始日の形式が不正です"))
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("終了日の形式が不正です"))
		return
	}

	var halfDayType *entity.HalfDayType
	if req.HalfDayType != nil {
		hdt := entity.HalfDayType(*req.HalfDayType)
		halfDayType = &hdt
	}

	input := &service.CreateLeaveRequestInput{
		EmployeeNumber: employeeNumber,
		LeaveType:      entity.LeaveType(req.LeaveType),
		StartDate:      startDate,
		EndDate:        endDate,
		HalfDayType:    halfDayType,
		Reason:         req.Reason,
	}

	leave, err := h.leaveService.CreateLeaveRequest(c.Request.Context(), input)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toLeaveRequestResponse(leave))
}

// GetLeaveRequests 休暇申請一覧取得
func (h *LeaveHandler) GetLeaveRequests(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")
	role := c.GetString("role")

	// 管理者の場合は他の従業員の申請も取得可能
	if role == "admin" {
		if empNum := c.Query("employee_number"); empNum != "" {
			num, err := strconv.Atoi(empNum)
			if err == nil {
				employeeNumber = num
			}
		}
	}

	var status *entity.ApprovalStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s, err := strconv.Atoi(statusStr)
		if err == nil {
			st := entity.ApprovalStatus(s)
			status = &st
		}
	}

	leaves, err := h.leaveService.GetLeaveRequests(c.Request.Context(), employeeNumber, status)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	responses := make([]LeaveRequestResponse, len(leaves))
	for i, leave := range leaves {
		responses[i] = toLeaveRequestResponse(leave)
	}

	c.JSON(http.StatusOK, responses)
}

// GetRemainingPaidLeaveDays 有給休暇残日数取得
func (h *LeaveHandler) GetRemainingPaidLeaveDays(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	remaining, err := h.leaveService.GetRemainingPaidLeaveDays(c.Request.Context(), employeeNumber)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"employee_number": employeeNumber,
		"remaining_days":  remaining,
		"total_days":      20.0, // TODO: 従業員マスタから取得
	})
}

func toLeaveRequestResponse(leave *entity.LeaveRequest) LeaveRequestResponse {
	resp := LeaveRequestResponse{
		ID:              leave.ID,
		EmployeeNumber:  leave.EmployeeNumber,
		LeaveType:       int(leave.LeaveType),
		StartDate:       leave.StartDate.Format("2006-01-02"),
		EndDate:         leave.EndDate.Format("2006-01-02"),
		Reason:          leave.Reason,
		ApprovalStatus:  int(leave.ApprovalStatus),
		ApproverID:      leave.ApproverID,
		ApprovalComment: leave.ApprovalComment,
		RejectReason:    leave.RejectReason,
		CreatedAt:       leave.CreatedAt,
	}

	if leave.HalfDayType != nil {
		hdt := int(*leave.HalfDayType)
		resp.HalfDayType = &hdt
	}

	if leave.ApprovedAt != nil {
		approvedAt := leave.ApprovedAt.Format(time.RFC3339)
		resp.ApprovedAt = &approvedAt
	}

	return resp
}
