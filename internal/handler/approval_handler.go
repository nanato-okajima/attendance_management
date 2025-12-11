package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
	"github.com/nanato-okajima/attendance_management/internal/service"
)

// ApprovalHandler 承認ハンドラー
type ApprovalHandler struct {
	approvalService   service.ApprovalService
	correctionService service.AttendanceCorrectionService
}

// NewApprovalHandler 承認ハンドラーを作成
func NewApprovalHandler(approvalService service.ApprovalService, correctionService service.AttendanceCorrectionService) *ApprovalHandler {
	return &ApprovalHandler{
		approvalService:   approvalService,
		correctionService: correctionService,
	}
}

// ApproveLeaveRequestRequest 承認リクエスト
type ApproveLeaveRequestRequest struct {
	Comment string `json:"comment"`
}

// RejectLeaveRequestRequest 却下リクエスト
type RejectLeaveRequestRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ApproveLeaveRequest 休暇申請承認
func (h *ApprovalHandler) ApproveLeaveRequest(c *gin.Context) {
	approverID := c.GetInt("employee_number")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("無効なIDです"))
		return
	}

	var req ApproveLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := h.approvalService.ApproveLeaveRequest(c.Request.Context(), uint(id), approverID, req.Comment); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "休暇申請を承認しました",
	})
}

// RejectLeaveRequest 休暇申請却下
func (h *ApprovalHandler) RejectLeaveRequest(c *gin.Context) {
	approverID := c.GetInt("employee_number")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("無効なIDです"))
		return
	}

	var req RejectLeaveRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := h.approvalService.RejectLeaveRequest(c.Request.Context(), uint(id), approverID, req.Reason); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "休暇申請を却下しました",
	})
}

// ApproveCorrectionRequest 打刻修正申請承認
func (h *ApprovalHandler) ApproveCorrectionRequest(c *gin.Context) {
	approverID := c.GetInt("employee_number")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("無効なIDです"))
		return
	}

	var req ApproveLeaveRequestRequest // 共通のリクエスト構造体を使用
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := h.correctionService.ApproveCorrectionRequest(c.Request.Context(), uint(id), approverID, req.Comment); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "打刻修正申請を承認しました",
	})
}

// RejectCorrectionRequest 打刻修正申請却下
func (h *ApprovalHandler) RejectCorrectionRequest(c *gin.Context) {
	approverID := c.GetInt("employee_number")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("無効なIDです"))
		return
	}

	var req RejectLeaveRequestRequest // 共通のリクエスト構造体を使用
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.NewValidationError(err.Error()))
		return
	}

	if err := h.correctionService.RejectCorrectionRequest(c.Request.Context(), uint(id), approverID, req.Reason); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "打刻修正申請を却下しました",
	})
}

// GetPendingApprovals 承認待ち一覧取得
func (h *ApprovalHandler) GetPendingApprovals(c *gin.Context) {
	leaves, err := h.approvalService.GetPendingApprovals(c.Request.Context())
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 承認待ちの打刻修正申請を取得
	pendingStatus := entity.ApprovalStatusPending
	corrections, err := h.correctionService.GetCorrectionRequests(c.Request.Context(), 0, &pendingStatus)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	leaveResponses := make([]LeaveRequestResponse, len(leaves))
	for i, leave := range leaves {
		leaveResponses[i] = toLeaveRequestResponse(leave)
	}

	correctionResponses := make([]CorrectionRequestResponse, len(corrections))
	for i, correction := range corrections {
		correctionResponses[i] = toCorrectionRequestResponse(correction)
	}

	c.JSON(http.StatusOK, gin.H{
		"leaves":      leaveResponses,
		"corrections": correctionResponses,
	})
}
