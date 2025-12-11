package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/nanato-okajima/attendance_management/internal/domain/attendance"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
)

type correctionRepository struct {
	db *gorm.DB
}

// NewCorrectionRepository 打刻修正リポジトリを作成
func NewCorrectionRepository(db *gorm.DB) attendance.CorrectionRepository {
	return &correctionRepository{db: db}
}

func (r *correctionRepository) Create(ctx context.Context, correction *entity.AttendanceCorrection) error {
	if err := r.db.WithContext(ctx).Create(correction).Error; err != nil {
		return errors.NewDBError(errors.DBErrorCreate, err)
	}
	return nil
}

func (r *correctionRepository) FindByID(ctx context.Context, id uint) (*entity.AttendanceCorrection, error) {
	var correction entity.AttendanceCorrection
	if err := r.db.WithContext(ctx).First(&correction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewDBError(errors.DBErrorNotFound, err)
		}
		return nil, errors.NewDBError(errors.DBErrorQuery, err)
	}
	return &correction, nil
}

func (r *correctionRepository) FindByEmployeeNumber(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.AttendanceCorrection, error) {
	var corrections []*entity.AttendanceCorrection
	query := r.db.WithContext(ctx).Where("employee_number = ?", employeeNumber)

	if status != nil {
		query = query.Where("approval_status = ?", *status)
	}

	if err := query.Order("created_at DESC").Find(&corrections).Error; err != nil {
		return nil, errors.NewDBError(errors.DBErrorQuery, err)
	}
	return corrections, nil
}

func (r *correctionRepository) FindPendingApprovals(ctx context.Context) ([]*entity.AttendanceCorrection, error) {
	var corrections []*entity.AttendanceCorrection
	if err := r.db.WithContext(ctx).
		Where("approval_status = ?", entity.ApprovalStatusPending).
		Order("created_at ASC").
		Find(&corrections).Error; err != nil {
		return nil, errors.NewDBError(errors.DBErrorQuery, err)
	}
	return corrections, nil
}

func (r *correctionRepository) Update(ctx context.Context, correction *entity.AttendanceCorrection) error {
	if err := r.db.WithContext(ctx).Save(correction).Error; err != nil {
		return errors.NewDBError(errors.DBErrorUpdate, err)
	}
	return nil
}
