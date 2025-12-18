package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/mock"
	"github.com/nanato-okajima/attendance_management/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLeaveService_CreateLeaveRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockLeaveRepository(ctrl)
	s := service.NewLeaveService(mockRepo)

	ctx := context.Background()
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 0, 1)

	tests := []struct {
		name    string
		input   *service.CreateLeaveRequestInput
		setup   func(m *mock.MockLeaveRepository)
		wantErr bool
	}{
		{
			name: "Success",
			input: &service.CreateLeaveRequestInput{
				EmployeeNumber: 1001,
				LeaveType:      entity.LeaveTypePaidLeave,
				StartDate:      startDate,
				EndDate:        endDate,
				Reason:         "Vacation",
			},
			setup: func(m *mock.MockLeaveRepository) {
				// CheckOverlap
				m.EXPECT().CheckOverlap(ctx, 1001, startDate, endDate, nil).Return(false, nil)
				// CountApprovedDays (for remaining days check)
				// Assuming annual allowance is 20.0, used 0.0
				m.EXPECT().CountApprovedDays(ctx, 1001, now.Year()).Return(0.0, nil)
				// Create
				m.EXPECT().Create(ctx, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ValidationError_StartDateAfterEndDate",
			input: &service.CreateLeaveRequestInput{
				EmployeeNumber: 1001,
				LeaveType:      entity.LeaveTypePaidLeave,
				StartDate:      endDate,
				EndDate:        startDate,
				Reason:         "Invalid dates",
			},
			setup: func(m *mock.MockLeaveRepository) {
				// No calls expected
			},
			wantErr: true,
		},
		{
			name: "OverlapError",
			input: &service.CreateLeaveRequestInput{
				EmployeeNumber: 1001,
				LeaveType:      entity.LeaveTypePaidLeave,
				StartDate:      startDate,
				EndDate:        endDate,
				Reason:         "Overlap",
			},
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().CheckOverlap(ctx, 1001, startDate, endDate, nil).Return(true, nil)
			},
			wantErr: true,
		},
		{
			name: "InsufficientPaidLeaveError",
			input: &service.CreateLeaveRequestInput{
				EmployeeNumber: 1001,
				LeaveType:      entity.LeaveTypePaidLeave,
				StartDate:      startDate,
				EndDate:        startDate.AddDate(0, 0, 20), // 21 days
				Reason:         "Too long",
			},
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().CheckOverlap(ctx, 1001, gomock.Any(), gomock.Any(), nil).Return(false, nil)
				// Used 0 days, allowance 20 days. Requesting 21 days.
				m.EXPECT().CountApprovedDays(ctx, 1001, now.Year()).Return(0.0, nil)
			},
			wantErr: true,
		},
		{
			name: "RepositoryError_Create",
			input: &service.CreateLeaveRequestInput{
				EmployeeNumber: 1001,
				LeaveType:      entity.LeaveTypeSpecialLeave,
				StartDate:      startDate,
				EndDate:        endDate,
				Reason:         "Special",
			},
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().CheckOverlap(ctx, 1001, startDate, endDate, nil).Return(false, nil)
				// Special leave doesn't check remaining days
				m.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockRepo)
			got, err := s.CreateLeaveRequest(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("LeaveService.CreateLeaveRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Equal(t, tt.input.EmployeeNumber, got.EmployeeNumber)
				assert.Equal(t, tt.input.LeaveType, got.LeaveType)
			}
		})
	}
}

func TestLeaveService_GetLeaveRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockLeaveRepository(ctrl)
	s := service.NewLeaveService(mockRepo)

	ctx := context.Background()
	employeeNumber := 1001
	status := entity.ApprovalStatusPending

	tests := []struct {
		name    string
		status  *entity.ApprovalStatus
		setup   func(m *mock.MockLeaveRepository)
		wantLen int
		wantErr bool
	}{
		{
			name:   "Success_WithStatus",
			status: &status,
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().FindByEmployeeNumber(ctx, employeeNumber, &status).Return([]*entity.LeaveRequest{
					{ID: 1, EmployeeNumber: employeeNumber, ApprovalStatus: status},
				}, nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "Success_NoStatus",
			status: nil,
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().FindByEmployeeNumber(ctx, employeeNumber, nil).Return([]*entity.LeaveRequest{
					{ID: 1, EmployeeNumber: employeeNumber},
					{ID: 2, EmployeeNumber: employeeNumber},
				}, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:   "RepositoryError",
			status: nil,
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().FindByEmployeeNumber(ctx, employeeNumber, nil).Return(nil, errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockRepo)
			got, err := s.GetLeaveRequests(ctx, employeeNumber, tt.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("LeaveService.GetLeaveRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Len(t, got, tt.wantLen)
			}
		})
	}
}

func TestLeaveService_GetRemainingPaidLeaveDays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockLeaveRepository(ctrl)
	s := service.NewLeaveService(mockRepo)

	ctx := context.Background()
	employeeNumber := 1001
	currentYear := time.Now().Year()

	tests := []struct {
		name    string
		setup   func(m *mock.MockLeaveRepository)
		want    float64
		wantErr bool
	}{
		{
			name: "Success",
			setup: func(m *mock.MockLeaveRepository) {
				// Used 5 days
				m.EXPECT().CountApprovedDays(ctx, employeeNumber, currentYear).Return(5.0, nil)
			},
			want:    15.0, // 20 - 5
			wantErr: false,
		},
		{
			name: "RepositoryError",
			setup: func(m *mock.MockLeaveRepository) {
				m.EXPECT().CountApprovedDays(ctx, employeeNumber, currentYear).Return(0.0, errors.New("db error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(mockRepo)
			got, err := s.GetRemainingPaidLeaveDays(ctx, employeeNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("LeaveService.GetRemainingPaidLeaveDays() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
