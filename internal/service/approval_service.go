package service

import (
	"context"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/domain/leave"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
)

// ApprovalService 承認サービスインターフェース
type ApprovalService interface {
	// ApproveLeaveRequest 休暇申請を承認
	ApproveLeaveRequest(ctx context.Context, id uint, approverID int, comment string) error

	// RejectLeaveRequest 休暇申請を却下
	RejectLeaveRequest(ctx context.Context, id uint, approverID int, reason string) error

	// GetPendingApprovals 承認待ち一覧を取得
	GetPendingApprovals(ctx context.Context) ([]*entity.LeaveRequest, error)
}

type approvalService struct {
	leaveRepo           leave.Repository
	notificationService NotificationService
}

// NewApprovalService 承認サービスを作成
func NewApprovalService(leaveRepo leave.Repository, notificationService NotificationService) ApprovalService {
	return &approvalService{
		leaveRepo:           leaveRepo,
		notificationService: notificationService,
	}
}

func (s *approvalService) ApproveLeaveRequest(ctx context.Context, id uint, approverID int, comment string) error {
	leaveReq, err := s.leaveRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if leaveReq.ApprovalStatus != entity.ApprovalStatusPending {
		return errors.NewValidationError("この申請は既に処理されています")
	}

	now := time.Now()
	leaveReq.ApprovalStatus = entity.ApprovalStatusApproved
	leaveReq.ApproverID = &approverID
	leaveReq.ApprovedAt = &now
	leaveReq.ApprovalComment = comment

	if err := s.leaveRepo.Update(ctx, leaveReq); err != nil {
		return err
	}

	// 通知送信
	s.notificationService.NotifyApprovalResult(ctx, leaveReq.EmployeeNumber, "承認", comment)

	return nil
}

func (s *approvalService) RejectLeaveRequest(ctx context.Context, id uint, approverID int, reason string) error {
	if reason == "" {
		return errors.NewValidationError("却下理由は必須です")
	}

	leaveReq, err := s.leaveRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if leaveReq.ApprovalStatus != entity.ApprovalStatusPending {
		return errors.NewValidationError("この申請は既に処理されています")
	}

	now := time.Now()
	leaveReq.ApprovalStatus = entity.ApprovalStatusRejected
	leaveReq.ApproverID = &approverID
	leaveReq.ApprovedAt = &now
	leaveReq.RejectReason = reason

	if err := s.leaveRepo.Update(ctx, leaveReq); err != nil {
		return err
	}

	// 通知送信
	s.notificationService.NotifyApprovalResult(ctx, leaveReq.EmployeeNumber, "却下", reason)

	return nil
}

func (s *approvalService) GetPendingApprovals(ctx context.Context) ([]*entity.LeaveRequest, error) {
	return s.leaveRepo.FindPendingApprovals(ctx)
}
