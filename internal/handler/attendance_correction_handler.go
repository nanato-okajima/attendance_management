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

// AttendanceCorrectionHandler 打刻修正ハンドラー
type AttendanceCorrectionHandler struct {
	correctionService service.AttendanceCorrectionService
}

// NewAttendanceCorrectionHandler 打刻修正ハンドラーを作成
func NewAttendanceCorrectionHandler(correctionService service.AttendanceCorrectionService) *AttendanceCorrectionHandler {
	return &AttendanceCorrectionHandler{
		correctionService: correctionService,
	}
}

// CreateCorrectionRequestRequest 打刻修正申請作成リクエスト
type CreateCorrectionRequestRequest struct {
	AttendanceID         uint   `json:"attendance_id" binding:"required"`
	CorrectionType       int    `json:"correction_type" binding:"required,min=1,max=3"`
	CorrectedOpeningTime string `json:"corrected_opening_time,omitempty"`
	CorrectedClosingTime string `json:"corrected_closing_time,omitempty"`
	Reason               string `json:"reason" binding:"required,max=500"`
}

// CorrectionRequestResponse 打刻修正申請レスポンス
type CorrectionRequestResponse struct {
	ID                   uint      `json:"id"`
	AttendanceID         uint      `json:"attendance_id"`
	EmployeeNumber       int       `json:"employee_number"`
	CorrectionType       int       `json:"correction_type"`
	CorrectedOpeningTime *string   `json:"corrected_opening_time,omitempty"`
	CorrectedClosingTime *string   `json:"corrected_closing_time,omitempty"`
	Reason               string    `json:"reason"`
	ApprovalStatus       int       `json:"approval_status"`
	ApproverID           *int      `json:"approver_id,omitempty"`
	ApprovedAt           *string   `json:"approved_at,omitempty"`
	RejectReason         string    `json:"reject_reason,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// CreateCorrectionRequest 打刻修正申請作成
func (h *AttendanceCorrectionHandler) CreateCorrectionRequest(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")

	var req CreateCorrectionRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	var openingTime, closingTime *time.Time
	if req.CorrectedOpeningTime != "" {
		t, err := time.Parse(time.RFC3339, req.CorrectedOpeningTime)
		if err != nil {
			errors.HandleError(c, errors.NewValidationError("出勤時刻の形式が不正です"))
			return
		}
		openingTime = &t
	}
	if req.CorrectedClosingTime != "" {
		t, err := time.Parse(time.RFC3339, req.CorrectedClosingTime)
		if err != nil {
			errors.HandleError(c, errors.NewValidationError("退勤時刻の形式が不正です"))
			return
		}
		closingTime = &t
	}

	input := &service.CreateCorrectionRequestInput{
		EmployeeNumber:       employeeNumber,
		AttendanceID:         req.AttendanceID,
		CorrectionType:       entity.CorrectionType(req.CorrectionType),
		CorrectedOpeningTime: openingTime,
		CorrectedClosingTime: closingTime,
		Reason:               req.Reason,
	}

	correction, err := h.correctionService.CreateCorrectionRequest(c.Request.Context(), input)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toCorrectionRequestResponse(correction))
}

// GetCorrectionRequests 打刻修正申請一覧取得
func (h *AttendanceCorrectionHandler) GetCorrectionRequests(c *gin.Context) {
	employeeNumber := c.GetInt("employee_number")
	role := c.GetString("role") // 管理者の場合は他の従業員の申請も取得可能
	if role == RoleAdmin {
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

	corrections, err := h.correctionService.GetCorrectionRequests(c.Request.Context(), employeeNumber, status)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	responses := make([]CorrectionRequestResponse, len(corrections))
	for i, c := range corrections {
		responses[i] = toCorrectionRequestResponse(c)
	}

	c.JSON(http.StatusOK, responses)
}

const RoleAdmin = "admin"

// GetCorrectionRequest 打刻修正申請詳細取得
func (h *AttendanceCorrectionHandler) GetCorrectionRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("無効なIDです"))
		return
	}

	correction, err := h.correctionService.GetCorrectionRequestByID(c.Request.Context(), uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 権限チェック: 本人または管理者のみ
	employeeNumber := c.GetInt("employee_number")
	role := c.GetString("role")
	if role != RoleAdmin && correction.EmployeeNumber != employeeNumber {
		errors.HandleError(c, errors.NewValidationError("権限がありません"))
		return
	}

	c.JSON(http.StatusOK, toCorrectionRequestResponse(correction))
}

func toCorrectionRequestResponse(c *entity.AttendanceCorrection) CorrectionRequestResponse {
	resp := CorrectionRequestResponse{
		ID:             c.ID,
		AttendanceID:   c.AttendanceID,
		EmployeeNumber: c.EmployeeNumber,
		CorrectionType: int(c.CorrectionType),
		Reason:         c.Reason,
		ApprovalStatus: int(c.ApprovalStatus),
		ApproverID:     c.ApproverID,
		RejectReason:   c.RejectReason,
		CreatedAt:      c.CreatedAt,
	}

	if c.CorrectedOpeningTime != nil {
		t := c.CorrectedOpeningTime.Format(time.RFC3339)
		resp.CorrectedOpeningTime = &t
	}
	if c.CorrectedClosingTime != nil {
		t := c.CorrectedClosingTime.Format(time.RFC3339)
		resp.CorrectedClosingTime = &t
	}
	if c.ApprovedAt != nil {
		t := c.ApprovedAt.Format(time.RFC3339)
		resp.ApprovedAt = &t
	}

	return resp
}
