package service

import (
	"context"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/domain/leave"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
)

//go:generate mockgen -source=leave_service.go -destination=../mock/mock_leave_service.go -package=mock

// LeaveService 休暇申請サービスインターフェース
type LeaveService interface {
	// CreateLeaveRequest 休暇申請を作成
	CreateLeaveRequest(ctx context.Context, req *CreateLeaveRequestInput) (*entity.LeaveRequest, error)

	// GetLeaveRequests 休暇申請一覧を取得
	GetLeaveRequests(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.LeaveRequest, error)

	// GetLeaveRequestByID 休暇申請詳細を取得
	GetLeaveRequestByID(ctx context.Context, id uint) (*entity.LeaveRequest, error)

	// GetRemainingPaidLeaveDays 有給休暇残日数を取得
	GetRemainingPaidLeaveDays(ctx context.Context, employeeNumber int) (float64, error)
}

// CreateLeaveRequestInput 休暇申請作成入力
type CreateLeaveRequestInput struct {
	EmployeeNumber int
	LeaveType      entity.LeaveType
	StartDate      time.Time
	EndDate        time.Time
	HalfDayType    *entity.HalfDayType
	Reason         string
}

type leaveService struct {
	leaveRepo leave.Repository
}

// NewLeaveService 休暇サービスを作成
func NewLeaveService(leaveRepo leave.Repository) LeaveService {
	return &leaveService{
		leaveRepo: leaveRepo,
	}
}

func (s *leaveService) CreateLeaveRequest(ctx context.Context, req *CreateLeaveRequestInput) (*entity.LeaveRequest, error) {
	// バリデーション: 開始日 <= 終了日
	if req.StartDate.After(req.EndDate) {
		return nil, errors.NewValidationError("開始日は終了日より前である必要があります")
	}

	// 期間の重複チェック
	hasOverlap, err := s.leaveRepo.CheckOverlap(ctx, req.EmployeeNumber, req.StartDate, req.EndDate, nil)
	if err != nil {
		return nil, err
	}
	if hasOverlap {
		return nil, errors.NewValidationError("申請期間が既存の申請と重複しています")
	}

	// 有給休暇の場合、残日数チェック
	if req.LeaveType == entity.LeaveTypePaidLeave {
		remaining, err := s.GetRemainingPaidLeaveDays(ctx, req.EmployeeNumber)
		if err != nil {
			return nil, err
		}

		requestedDays := calculateDays(req.StartDate, req.EndDate, req.HalfDayType)
		if requestedDays > remaining {
			return nil, errors.NewValidationError("有給休暇の残日数が不足しています")
		}
	}

	// 休暇申請作成
	leaveRequest := &entity.LeaveRequest{
		EmployeeNumber: req.EmployeeNumber,
		LeaveType:      req.LeaveType,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		HalfDayType:    req.HalfDayType,
		Reason:         req.Reason,
		ApprovalStatus: entity.ApprovalStatusPending,
	}

	if err := s.leaveRepo.Create(ctx, leaveRequest); err != nil {
		return nil, err
	}

	return leaveRequest, nil
}

func (s *leaveService) GetLeaveRequests(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.LeaveRequest, error) {
	return s.leaveRepo.FindByEmployeeNumber(ctx, employeeNumber, status)
}

func (s *leaveService) GetLeaveRequestByID(ctx context.Context, id uint) (*entity.LeaveRequest, error) {
	return s.leaveRepo.FindByID(ctx, id)
}

func (s *leaveService) GetRemainingPaidLeaveDays(ctx context.Context, employeeNumber int) (float64, error) {
	// TODO: 実際の付与日数は従業員マスタから取得
	const annualAllowance = 20.0

	currentYear := time.Now().Year()
	usedDays, err := s.leaveRepo.CountApprovedDays(ctx, employeeNumber, currentYear)
	if err != nil {
		return 0, err
	}

	return annualAllowance - usedDays, nil
}

// calculateDays 休暇日数を計算
func calculateDays(startDate, endDate time.Time, halfDayType *entity.HalfDayType) float64 {
	days := endDate.Sub(startDate).Hours()/24 + 1

	if halfDayType != nil {
		return days - 0.5
	}

	return days
}
