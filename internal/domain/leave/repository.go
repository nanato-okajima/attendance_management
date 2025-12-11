package leave

import (
	"context"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/entity"
)

//go:generate mockgen -source=repository.go -destination=../../mock/mock_leave_repository.go -package=mock -mock_names Repository=MockLeaveRepository

// Repository 休暇申請リポジトリインターフェース
type Repository interface {
	// Create 休暇申請を作成
	Create(ctx context.Context, leave *entity.LeaveRequest) error

	// FindByID IDで休暇申請を取得
	FindByID(ctx context.Context, id uint) (*entity.LeaveRequest, error)

	// FindByEmployeeNumber 従業員番号で休暇申請一覧を取得
	FindByEmployeeNumber(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.LeaveRequest, error)

	// FindPendingApprovals 承認待ちの申請一覧を取得
	FindPendingApprovals(ctx context.Context) ([]*entity.LeaveRequest, error)

	// Update 休暇申請を更新
	Update(ctx context.Context, leave *entity.LeaveRequest) error

	// CheckOverlap 期間の重複をチェック
	CheckOverlap(ctx context.Context, employeeNumber int, startDate, endDate time.Time, excludeID *uint) (bool, error)

	// CountApprovedDays 承認済みの有給休暇日数を取得
	CountApprovedDays(ctx context.Context, employeeNumber int, year int) (float64, error)
}
