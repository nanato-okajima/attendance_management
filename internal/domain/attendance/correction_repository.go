package attendance

import (
	"context"

	"github.com/nanato-okajima/attendance_management/internal/entity"
)

//go:generate mockgen -source=correction_repository.go -destination=../../mock/mock_correction_repository.go -package=mock -mock_names CorrectionRepository=MockCorrectionRepository

// CorrectionRepository 打刻修正リポジトリインターフェース
type CorrectionRepository interface {
	// Create 打刻修正申請を作成
	Create(ctx context.Context, correction *entity.AttendanceCorrection) error

	// FindByID IDで打刻修正申請を取得
	FindByID(ctx context.Context, id uint) (*entity.AttendanceCorrection, error)

	// FindByEmployeeNumber 従業員番号で打刻修正申請一覧を取得
	FindByEmployeeNumber(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.AttendanceCorrection, error)

	// FindPendingApprovals 承認待ちの申請一覧を取得
	FindPendingApprovals(ctx context.Context) ([]*entity.AttendanceCorrection, error)

	// Update 打刻修正申請を更新
	Update(ctx context.Context, correction *entity.AttendanceCorrection) error
}
