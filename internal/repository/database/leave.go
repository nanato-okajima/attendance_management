package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/nanato-okajima/attendance_management/internal/domain/leave"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
)

type leaveRepository struct {
	db *gorm.DB
}

// NewLeaveRepository 休暇リポジトリを作成
func NewLeaveRepository(db *gorm.DB) leave.Repository {
	return &leaveRepository{db: db}
}

func (r *leaveRepository) Create(ctx context.Context, leave *entity.LeaveRequest) error {
	if err := r.db.WithContext(ctx).Create(leave).Error; err != nil {
		return errors.NewDBError(errors.DBErrorCreate, err)
	}
	return nil
}

func (r *leaveRepository) FindByID(ctx context.Context, id uint) (*entity.LeaveRequest, error) {
	var leave entity.LeaveRequest
	if err := r.db.WithContext(ctx).First(&leave, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewDBError(errors.DBErrorNotFound, err)
		}
		return nil, errors.NewDBError(errors.DBErrorQuery, err)
	}
	return &leave, nil
}

func (r *leaveRepository) FindByEmployeeNumber(ctx context.Context, employeeNumber int, status *entity.ApprovalStatus) ([]*entity.LeaveRequest, error) {
	var leaves []*entity.LeaveRequest
	query := r.db.WithContext(ctx).Where("employee_number = ?", employeeNumber)

	if status != nil {
		query = query.Where("approval_status = ?", *status)
	}

	if err := query.Order("start_date DESC").Find(&leaves).Error; err != nil {
		return nil, errors.NewDBError(errors.DBErrorQuery, err)
	}
	return leaves, nil
}

func (r *leaveRepository) FindPendingApprovals(ctx context.Context) ([]*entity.LeaveRequest, error) {
	var leaves []*entity.LeaveRequest
	if err := r.db.WithContext(ctx).
		Where("approval_status = ?", entity.ApprovalStatusPending).
		Order("created_at ASC").
		Find(&leaves).Error; err != nil {
		return nil, errors.NewDBError(errors.DBErrorQuery, err)
	}
	return leaves, nil
}

func (r *leaveRepository) Update(ctx context.Context, leave *entity.LeaveRequest) error {
	if err := r.db.WithContext(ctx).Save(leave).Error; err != nil {
		return errors.NewDBError(errors.DBErrorUpdate, err)
	}
	return nil
}

func (r *leaveRepository) CheckOverlap(ctx context.Context, employeeNumber int, startDate, endDate time.Time, excludeID *uint) (bool, error) {
	query := r.db.WithContext(ctx).Model(&entity.LeaveRequest{}).
		Where("employee_number = ?", employeeNumber).
		Where("approval_status != ?", entity.ApprovalStatusRejected).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, errors.NewDBError(errors.DBErrorQuery, err)
	}

	return count > 0, nil
}

func (r *leaveRepository) CountApprovedDays(ctx context.Context, employeeNumber int, year int) (float64, error) {
	var result struct {
		TotalDays float64
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.Local)

	err := r.db.WithContext(ctx).Model(&entity.LeaveRequest{}).
		Select("COALESCE(SUM(DATEDIFF(end_date, start_date) + 1), 0) as total_days").
		Where("employee_number = ?", employeeNumber).
		Where("leave_type = ?", entity.LeaveTypePaidLeave).
		Where("approval_status = ?", entity.ApprovalStatusApproved).
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return 0, errors.NewDBError(errors.DBErrorQuery, err)
	}

	return result.TotalDays, nil
}
