package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/domain/attendance"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/mock"
	"github.com/nanato-okajima/attendance_management/internal/service"
)

func TestAttendanceService_ClockIn(t *testing.T) {
	t.Parallel()
	// 固定時刻を設定 (2025-12-01 09:00:00)
	fixedTime := time.Date(2025, 12, 1, 9, 0, 0, 0, time.Local)
	todayDate := time.Date(2025, 12, 1, 0, 0, 0, 0, time.Local)
	timeProvider := &service.FixedTimeProvider{FixedTime: fixedTime}
	latitude := 35.6812
	longitude := 139.7671

	testConfig := &config.WorkHoursConfig{
		StartHour:         9,
		StartMinute:       0,
		EndHour:           18,
		EndMinute:         0,
		BreakHours:        1,
		StandardWorkHours: 8,
	}

	type param struct {
		employeeID  int
		latitude    *float64
		longitude   *float64
		clockSource int
	}

	tests := []struct {
		name          string
		param         param
		setupMock     func(*mock.MockAttendanceRepository)
		hasError      bool
		expected      *entity.Attendance
		expectedError error
	}{
		{
			name: "新規出勤報告",
			param: param{
				employeeID:  1,
				latitude:    &latitude,
				longitude:   &longitude,
				clockSource: 1,
			},
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				mockRepo.EXPECT().FindByEmployeeAndDate(1, todayDate).Return(nil, gorm.ErrRecordNotFound)
				mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
			},
			expected: &entity.Attendance{
				EmployeeID:       1,
				TargetDate:       todayDate,
				OpeningTime:      lo.ToPtr(fixedTime),
				AttendanceStatus: int(attendance.StatusNormal),
				ClockSource:      1,
				Latitude:         &latitude,
				Longitude:        &longitude,
			},
			hasError: false,
		},
		{
			name: "既存レコード更新（OpeningTimeなし）",
			param: param{
				employeeID:  1,
				latitude:    &latitude,
				longitude:   &longitude,
				clockSource: 1,
			},
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				existingAttendance := &entity.Attendance{
					AttendanceID: 10,
					EmployeeID:   1,
					TargetDate:   todayDate,
					OpeningTime:  nil, // まだ出勤していない
				}
				mockRepo.EXPECT().FindByEmployeeAndDate(1, todayDate).Return(existingAttendance, nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expected: &entity.Attendance{
				AttendanceID:     10,
				EmployeeID:       1,
				TargetDate:       todayDate,
				OpeningTime:      lo.ToPtr(fixedTime),
				AttendanceStatus: int(attendance.StatusNormal),
				ClockSource:      1,
				Latitude:         &latitude,
				Longitude:        &longitude,
			},
			hasError: false,
		},
		{
			name: "既に出勤報告済み",
			param: param{
				employeeID:  1,
				latitude:    nil,
				longitude:   nil,
				clockSource: 1,
			},
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				existingAttendance := &entity.Attendance{
					AttendanceID: 1,
					EmployeeID:   1,
					OpeningTime:  &fixedTime,
				}
				mockRepo.EXPECT().FindByEmployeeAndDate(1, todayDate).Return(existingAttendance, nil)
			},
			hasError:      true,
			expectedError: errors.New("[ATT_001] Already clocked in today"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockAttendanceRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := service.NewAttendanceServiceWithTimeProvider(mockRepo, testConfig, timeProvider)

			got, err := svc.ClockIn(tt.param.employeeID, tt.param.latitude, tt.param.longitude, tt.param.clockSource)
			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, got)
				// if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				// }
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestAttendanceService_ClockOut(t *testing.T) {
	t.Parallel()
	// 固定時刻を設定 (2025-12-01 18:00:00)
	fixedTime := time.Date(2025, 12, 1, 18, 0, 0, 0, time.Local)
	timeProvider := &service.FixedTimeProvider{FixedTime: fixedTime}

	testConfig := &config.WorkHoursConfig{
		StartHour:         9,
		StartMinute:       0,
		EndHour:           18,
		EndMinute:         0,
		BreakHours:        1,
		StandardWorkHours: 8,
	}

	type param struct {
		employeeID int
		latitude   *float64
		longitude  *float64
	}

	tests := []struct {
		name          string
		param         param
		setupMock     func(*mock.MockAttendanceRepository)
		hasError      bool
		expectedError string
	}{
		{
			name: "Success",
			param: param{
				employeeID: 1,
				latitude:   nil,
				longitude:  nil,
			},
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				openingTime := time.Date(2025, 12, 1, 9, 0, 0, 0, time.Local)
				existingAttendance := &entity.Attendance{
					AttendanceID:     1,
					EmployeeID:       1,
					OpeningTime:      &openingTime,
					ClosingTime:      nil,
					AttendanceStatus: int(attendance.StatusNormal),
				}
				mockRepo.EXPECT().
					FindByEmployeeAndDate(1, time.Date(2025, 12, 1, 0, 0, 0, 0, time.Local)).
					Return(existingAttendance, nil)
				mockRepo.EXPECT().
					Update(gomock.Any()).
					Return(nil)
			},
			hasError: false,
		},
		{
			name: "出勤報告なし",
			param: param{
				employeeID: 1,
				latitude:   nil,
				longitude:  nil,
			},
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				mockRepo.EXPECT().
					FindByEmployeeAndDate(1, time.Date(2025, 12, 1, 0, 0, 0, 0, time.Local)).
					Return(nil, errors.New("not found"))
			},
			hasError:      true,
			expectedError: "[ATT_002] No clock-in record found",
		},
		{
			name: "既に退勤報告済み",
			param: param{
				employeeID: 1,
				latitude:   nil,
				longitude:  nil,
			},
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				openingTime := time.Date(2025, 12, 1, 9, 0, 0, 0, time.Local)
				closingTime := time.Date(2025, 12, 1, 17, 0, 0, 0, time.Local)
				existingAttendance := &entity.Attendance{
					AttendanceID: 1,
					EmployeeID:   1,
					OpeningTime:  &openingTime,
					ClosingTime:  &closingTime,
				}
				mockRepo.EXPECT().
					FindByEmployeeAndDate(1, time.Date(2025, 12, 1, 0, 0, 0, 0, time.Local)).
					Return(existingAttendance, nil)
			},
			hasError:      true,
			expectedError: "[ATT_003] Already clocked out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockAttendanceRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := service.NewAttendanceServiceWithTimeProvider(mockRepo, testConfig, timeProvider)

			attendance, err := svc.ClockOut(tt.param.employeeID, tt.param.latitude, tt.param.longitude)
			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, attendance)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, attendance)
				assert.NotNil(t, attendance.ClosingTime)
				assert.NotNil(t, attendance.WorkHours)
			}
		})
	}
}

func TestAttendanceService_GetTodayAttendance(t *testing.T) {
	t.Parallel()
	// 固定時刻を設定
	fixedTime := time.Date(2025, 12, 1, 9, 0, 0, 0, time.Local)
	todayDate := time.Date(2025, 12, 1, 0, 0, 0, 0, time.Local)
	timeProvider := &service.FixedTimeProvider{FixedTime: fixedTime}

	testConfig := &config.WorkHoursConfig{
		StartHour:         9,
		StartMinute:       0,
		EndHour:           18,
		EndMinute:         0,
		BreakHours:        1,
		StandardWorkHours: 8,
	}

	tests := []struct {
		name               string
		employeeID         int
		setupMock          func(*mock.MockAttendanceRepository)
		hasError           bool
		expectedAttendance *entity.Attendance
	}{
		{
			name:       "Success",
			employeeID: 1,
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				expectedAttendance := &entity.Attendance{
					AttendanceID: 1,
					EmployeeID:   1,
					OpeningTime:  &fixedTime,
				}
				mockRepo.EXPECT().FindByEmployeeAndDate(1, todayDate).Return(expectedAttendance, nil)
			},
			hasError: false,
			expectedAttendance: &entity.Attendance{
				AttendanceID: 1,
				EmployeeID:   1,
				OpeningTime:  &fixedTime,
			},
		},
		{
			name:       "NotFound",
			employeeID: 1,
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				mockRepo.EXPECT().FindByEmployeeAndDate(1, todayDate).Return(nil, gorm.ErrRecordNotFound)
			},
			hasError:           true,
			expectedAttendance: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockAttendanceRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := service.NewAttendanceServiceWithTimeProvider(mockRepo, testConfig, timeProvider)

			attendance, err := svc.GetTodayAttendance(tt.employeeID)
			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, attendance)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAttendance, attendance)
			}
		})
	}
}

func TestAttendanceService_GetMonthlyAttendances(t *testing.T) {
	t.Parallel()
	testConfig := &config.WorkHoursConfig{
		StartHour:         9,
		StartMinute:       0,
		EndHour:           18,
		EndMinute:         0,
		BreakHours:        1,
		StandardWorkHours: 8,
	}

	tests := []struct {
		name                string
		employeeID          int
		year                int
		month               int
		setupMock           func(*mock.MockAttendanceRepository)
		hasError            bool
		expectedCount       int
		expectedAttendances []entity.Attendance
	}{
		{
			name:       "Success",
			employeeID: 1,
			year:       2025,
			month:      11,
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				expectedAttendances := []entity.Attendance{
					{AttendanceID: 1, EmployeeID: 1},
					{AttendanceID: 2, EmployeeID: 1},
				}
				mockRepo.EXPECT().FindByEmployeeAndDateRange(1, gomock.Any(), gomock.Any()).Return(expectedAttendances, nil)
			},
			hasError:      false,
			expectedCount: 2,
		},
		{
			name:       "EmptyResult",
			employeeID: 1,
			year:       2025,
			month:      12,
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				mockRepo.EXPECT().FindByEmployeeAndDateRange(1, gomock.Any(), gomock.Any()).Return([]entity.Attendance{}, nil)
			},
			hasError:      false,
			expectedCount: 0,
		},
		{
			name:       "RepositoryError",
			employeeID: 1,
			year:       2025,
			month:      11,
			setupMock: func(mockRepo *mock.MockAttendanceRepository) {
				mockRepo.EXPECT().FindByEmployeeAndDateRange(1, gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			hasError:      true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockAttendanceRepository(ctrl)
			tt.setupMock(mockRepo)

			svc := service.NewAttendanceService(mockRepo, testConfig)

			attendances, err := svc.GetMonthlyAttendances(tt.employeeID, tt.year, tt.month)
			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, attendances)
			} else {
				assert.NoError(t, err)
				assert.Len(t, attendances, tt.expectedCount)
			}
		})
	}
}
