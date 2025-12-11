package service

import (
	"context"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/domain/attendance"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
)

// AttendanceCorrectionService 打刻修正サービスインターフェース
type AttendanceCorrectionService interface {
	// CreateCorrectionRequest 打刻修正申請を作成
	CreateCorrectionRequest(ctx context.Context, req *CreateCorrectionRequestInput) (*entity.AttendanceCorrection, error)

	// GetCorrectionRequests 打刻修正申請一覧を取得
	GetCorrectionRequests(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.AttendanceCorrection, error)

	// GetCorrectionRequestByID 打刻修正申請詳細を取得
	GetCorrectionRequestByID(ctx context.Context, id uint) (*entity.AttendanceCorrection, error)

	// ApproveCorrectionRequest 打刻修正申請を承認
	ApproveCorrectionRequest(ctx context.Context, id uint, approverID int, comment string) error

	// RejectCorrectionRequest 打刻修正申請を却下
	RejectCorrectionRequest(ctx context.Context, id uint, approverID int, reason string) error
}

// CreateCorrectionRequestInput 打刻修正申請作成入力
type CreateCorrectionRequestInput struct {
	EmployeeNumber       int
	AttendanceID         uint
	CorrectionType       entity.CorrectionType
	CorrectedOpeningTime *time.Time
	CorrectedClosingTime *time.Time
	Reason               string
}

type attendanceCorrectionService struct {
	correctionRepo      attendance.CorrectionRepository
	attendanceRepo      attendance.Repository
	notificationService NotificationService
}

// NewAttendanceCorrectionService 打刻修正サービスを作成
func NewAttendanceCorrectionService(
	correctionRepo attendance.CorrectionRepository,
	attendanceRepo attendance.Repository,
	notificationService NotificationService,
) AttendanceCorrectionService {
	return &attendanceCorrectionService{
		correctionRepo:      correctionRepo,
		attendanceRepo:      attendanceRepo,
		notificationService: notificationService,
	}
}

func (s *attendanceCorrectionService) CreateCorrectionRequest(ctx context.Context, req *CreateCorrectionRequestInput) (*entity.AttendanceCorrection, error) {
	// 対象の勤怠記録が存在するか確認
	// TODO: AttendanceRepositoryにFindByIDを追加する必要があるかもしれないが、
	// ここではAttendanceIDが有効であることを前提とするか、あるいは修正対象の日付から検索するロジックにするか。
	// 要件では「対象日」を入力としているが、実装計画ではAttendanceIDを受け取っている。
	// ここではシンプルに作成処理を行う。

	// バリデーション: 修正時刻の整合性
	if req.CorrectionType == entity.CorrectionTypeBoth {
		if req.CorrectedOpeningTime == nil || req.CorrectedClosingTime == nil {
			return nil, errors.NewValidationError("出退勤時刻の両方が必要です")
		}
		if req.CorrectedOpeningTime.After(*req.CorrectedClosingTime) {
			return nil, errors.NewValidationError("出勤時刻は退勤時刻より前である必要があります")
		}
	}

	correction := &entity.AttendanceCorrection{
		EmployeeNumber:       req.EmployeeNumber,
		AttendanceID:         req.AttendanceID,
		CorrectionType:       req.CorrectionType,
		CorrectedOpeningTime: req.CorrectedOpeningTime,
		CorrectedClosingTime: req.CorrectedClosingTime,
		Reason:               req.Reason,
		ApprovalStatus:       entity.ApprovalStatusPending,
	}

	if err := s.correctionRepo.Create(ctx, correction); err != nil {
		return nil, err
	}

	// 承認待ち通知
	// TODO: 承認者のIDを特定するロジックが必要（例：部署のマネージャー）
	// ここでは仮に0としている
	s.notificationService.NotifyPendingApproval(ctx, 0, "打刻修正", "新規申請")

	return correction, nil
}

func (s *attendanceCorrectionService) GetCorrectionRequests(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.AttendanceCorrection, error) {
	return s.correctionRepo.FindByEmployeeNumber(ctx, employeeNumber, status)
}

func (s *attendanceCorrectionService) GetCorrectionRequestByID(ctx context.Context, id uint) (*entity.AttendanceCorrection, error) {
	return s.correctionRepo.FindByID(ctx, id)
}

func (s *attendanceCorrectionService) ApproveCorrectionRequest(ctx context.Context, id uint, approverID int, comment string) error {
	correction, err := s.correctionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if correction.ApprovalStatus != entity.ApprovalStatusPending {
		return errors.NewValidationError("この申請は既に処理されています")
	}

	// 勤怠記録を更新
	attendance, err := s.attendanceRepo.FindByID(correction.AttendanceID)
	if err != nil {
		return err
	}

	if correction.CorrectionType == entity.CorrectionTypeClockIn || correction.CorrectionType == entity.CorrectionTypeBoth {
		if correction.CorrectedOpeningTime != nil {
			attendance.OpeningTime = correction.CorrectedOpeningTime
		}
	}
	if correction.CorrectionType == entity.CorrectionTypeClockOut || correction.CorrectionType == entity.CorrectionTypeBoth {
		if correction.CorrectedClosingTime != nil {
			attendance.ClosingTime = correction.CorrectedClosingTime
		}
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return err
	}

	now := time.Now()
	correction.ApprovalStatus = entity.ApprovalStatusApproved
	correction.ApproverID = &approverID
	correction.ApprovedAt = &now
	// correction.ApprovalComment = comment

	if err := s.correctionRepo.Update(ctx, correction); err != nil {
		return err
	}

	s.notificationService.NotifyApprovalResult(ctx, correction.EmployeeNumber, "承認", comment)

	return nil
}

func (s *attendanceCorrectionService) RejectCorrectionRequest(ctx context.Context, id uint, approverID int, reason string) error {
	if reason == "" {
		return errors.NewValidationError("却下理由は必須です")
	}

	correction, err := s.correctionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if correction.ApprovalStatus != entity.ApprovalStatusPending {
		return errors.NewValidationError("この申請は既に処理されています")
	}

	now := time.Now()
	correction.ApprovalStatus = entity.ApprovalStatusRejected
	correction.ApproverID = &approverID
	correction.ApprovedAt = &now
	correction.RejectReason = reason

	if err := s.correctionRepo.Update(ctx, correction); err != nil {
		return err
	}

	s.notificationService.NotifyApprovalResult(ctx, correction.EmployeeNumber, "却下", reason)

	return nil
}
