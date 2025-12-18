package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/mock"
	"github.com/nanato-okajima/attendance_management/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApprovalService_ApproveLeaveRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLeaveRepo := mock.NewMockLeaveRepository(ctrl)
	mockNotifyService := mock.NewMockNotificationService(ctrl)
	s := service.NewApprovalService(mockLeaveRepo, mockNotifyService)

	ctx := context.Background()
	id := uint(1)
	approverID := 9999
	comment := "Approved"

	tests := []struct {
		name    string
		setup   func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService)
		wantErr bool
	}{
		{
			name: "Success",
			setup: func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService) {
				leaveReq := &entity.LeaveRequest{
					ID:             id,
					EmployeeNumber: 1001,
					ApprovalStatus: entity.ApprovalStatusPending,
				}
				ml.EXPECT().FindByID(ctx, id).Return(leaveReq, nil)
				ml.EXPECT().Update(ctx, gomock.Any()).Return(nil)
				mn.EXPECT().NotifyApprovalResult(ctx, 1001, "承認", comment)
			},
			wantErr: false,
		},
		{
			name: "AlreadyProcessed",
			setup: func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService) {
				leaveReq := &entity.LeaveRequest{
					ID:             id,
					ApprovalStatus: entity.ApprovalStatusApproved,
				}
				ml.EXPECT().FindByID(ctx, id).Return(leaveReq, nil)
			},
			wantErr: true,
		},
		{
			name: "RepositoryError_FindByID",
			setup: func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService) {
				ml.EXPECT().FindByID(ctx, id).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockLeaveRepo, mockNotifyService)
			err := s.ApproveLeaveRequest(ctx, id, approverID, comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovalService.ApproveLeaveRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApprovalService_RejectLeaveRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLeaveRepo := mock.NewMockLeaveRepository(ctrl)
	mockNotifyService := mock.NewMockNotificationService(ctrl)
	s := service.NewApprovalService(mockLeaveRepo, mockNotifyService)

	ctx := context.Background()
	id := uint(1)
	approverID := 9999
	reason := "Rejected"

	tests := []struct {
		name    string
		reason  string
		setup   func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService)
		wantErr bool
	}{
		{
			name:   "Success",
			reason: reason,
			setup: func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService) {
				leaveReq := &entity.LeaveRequest{
					ID:             id,
					EmployeeNumber: 1001,
					ApprovalStatus: entity.ApprovalStatusPending,
				}
				ml.EXPECT().FindByID(ctx, id).Return(leaveReq, nil)
				ml.EXPECT().Update(ctx, gomock.Any()).Return(nil)
				mn.EXPECT().NotifyApprovalResult(ctx, 1001, "却下", reason)
			},
			wantErr: false,
		},
		{
			name:   "ValidationError_EmptyReason",
			reason: "",
			setup: func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService) {
				// No calls expected
			},
			wantErr: true,
		},
		{
			name:   "AlreadyProcessed",
			reason: reason,
			setup: func(ml *mock.MockLeaveRepository, mn *mock.MockNotificationService) {
				leaveReq := &entity.LeaveRequest{
					ID:             id,
					ApprovalStatus: entity.ApprovalStatusRejected,
				}
				ml.EXPECT().FindByID(ctx, id).Return(leaveReq, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockLeaveRepo, mockNotifyService)
			err := s.RejectLeaveRequest(ctx, id, approverID, tt.reason)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovalService.RejectLeaveRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApprovalService_GetPendingApprovals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLeaveRepo := mock.NewMockLeaveRepository(ctrl)
	mockNotifyService := mock.NewMockNotificationService(ctrl)
	s := service.NewApprovalService(mockLeaveRepo, mockNotifyService)

	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(ml *mock.MockLeaveRepository)
		wantLen int
		wantErr bool
	}{
		{
			name: "Success",
			setup: func(ml *mock.MockLeaveRepository) {
				ml.EXPECT().FindPendingApprovals(ctx).Return([]*entity.LeaveRequest{
					{ID: 1}, {ID: 2},
				}, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "RepositoryError",
			setup: func(ml *mock.MockLeaveRepository) {
				ml.EXPECT().FindPendingApprovals(ctx).Return(nil, errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockLeaveRepo)
			got, err := s.GetPendingApprovals(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApprovalService.GetPendingApprovals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Len(t, got, tt.wantLen)
			}
		})
	}
}
