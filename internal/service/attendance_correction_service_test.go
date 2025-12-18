package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/mock"
	"github.com/nanato-okajima/attendance_management/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAttendanceCorrectionService_CreateCorrectionRequest(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	attendanceID := uint(1)
	employeeNumber := 1001
	reason := "Forgot to clock out"
	now := time.Now()

	tests := []struct {
		name    string
		input   *service.CreateCorrectionRequestInput
		setup   func(mc *mock.MockCorrectionRepository, ma *mock.MockAttendanceRepository, mn *mock.MockNotificationService)
		wantErr bool
	}{
		{
			name: "Success",
			input: &service.CreateCorrectionRequestInput{
				AttendanceID:         attendanceID,
				EmployeeNumber:       employeeNumber,
				CorrectionType:       entity.CorrectionTypeClockOut,
				CorrectedClosingTime: &now,
				Reason:               reason,
			},
			setup: func(mc *mock.MockCorrectionRepository, ma *mock.MockAttendanceRepository, mn *mock.MockNotificationService) {
				mc.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mn.EXPECT().NotifyPendingApproval(ctx, 0, "打刻修正", "新規申請")
			},
			wantErr: false,
		},
		/*
			{
				name: "AttendanceNotFound",
				input: &service.CreateCorrectionRequestInput{
					AttendanceID:   attendanceID,
					EmployeeNumber: employeeNumber,
				},
				setup: func(mc *mock.MockCorrectionRepository, ma *mock.MockAttendanceRepository, mn *mock.MockNotificationService) {
					ma.EXPECT().FindByID(attendanceID).Return(nil, errors.New("not found"))
				},
				wantErr: true,
			},
			{
				name: "Unauthorized",
				input: &service.CreateCorrectionRequestInput{
					AttendanceID:   attendanceID,
					EmployeeNumber: 1002, // Different employee
				},
				setup: func(mc *mock.MockCorrectionRepository, ma *mock.MockAttendanceRepository, mn *mock.MockNotificationService) {
					ma.EXPECT().FindByID(attendanceID).Return(&entity.Attendance{AttendanceID: attendanceID, EmployeeID: employeeNumber}, nil)
				},
				wantErr: true,
			},
		*/
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCorrectionRepo := mock.NewMockCorrectionRepository(ctrl)
			mockAttendanceRepo := mock.NewMockAttendanceRepository(ctrl)
			mockNotifyService := mock.NewMockNotificationService(ctrl)
			s := service.NewAttendanceCorrectionService(mockCorrectionRepo, mockAttendanceRepo, mockNotifyService)

			tt.setup(mockCorrectionRepo, mockAttendanceRepo, mockNotifyService)
			_, err := s.CreateCorrectionRequest(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttendanceCorrectionService.CreateCorrectionRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttendanceCorrectionService_ApproveCorrectionRequest(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	id := uint(1)
	approverID := 9999
	comment := "Approved"
	attendanceID := uint(100)
	now := time.Now()

	tests := []struct {
		name    string
		setup   func(mc *mock.MockCorrectionRepository, ma *mock.MockAttendanceRepository, mn *mock.MockNotificationService)
		wantErr bool
	}{
		{
			name: "Success_ClockOut",
			setup: func(mc *mock.MockCorrectionRepository, ma *mock.MockAttendanceRepository, mn *mock.MockNotificationService) {
				correction := &entity.AttendanceCorrection{
					ID:                   id,
					AttendanceID:         attendanceID,
					EmployeeNumber:       1001,
					ApprovalStatus:       entity.ApprovalStatusPending,
					CorrectionType:       entity.CorrectionTypeClockOut,
					CorrectedClosingTime: &now,
				}
				mc.EXPECT().FindByID(ctx, id).Return(correction, nil)

				attendance := &entity.Attendance{AttendanceID: attendanceID}
				ma.EXPECT().FindByID(attendanceID).Return(attendance, nil)

				// Expect Update on Attendance
				ma.EXPECT().Update(gomock.Any()).DoAndReturn(func(a *entity.Attendance) error {
					assert.Equal(t, now, *a.ClosingTime)
					return nil
				})

				// Expect Update on Correction
				mc.EXPECT().Update(ctx, gomock.Any()).Return(nil)

				mn.EXPECT().NotifyApprovalResult(ctx, 1001, "承認", comment)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCorrectionRepo := mock.NewMockCorrectionRepository(ctrl)
			mockAttendanceRepo := mock.NewMockAttendanceRepository(ctrl)
			mockNotifyService := mock.NewMockNotificationService(ctrl)
			s := service.NewAttendanceCorrectionService(mockCorrectionRepo, mockAttendanceRepo, mockNotifyService)

			tt.setup(mockCorrectionRepo, mockAttendanceRepo, mockNotifyService)
			err := s.ApproveCorrectionRequest(ctx, id, approverID, comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("AttendanceCorrectionService.ApproveCorrectionRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
