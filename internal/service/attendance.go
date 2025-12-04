package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/nanato-okajima/attendance_management/internal/config"
	"github.com/nanato-okajima/attendance_management/internal/domain/attendance"
	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
)

type AttendanceService interface {
	ClockIn(employeeID int, latitude, longitude *float64, clockSource int) (*entity.Attendance, error)
	ClockOut(employeeID int, latitude, longitude *float64) (*entity.Attendance, error)
	GetTodayAttendance(employeeID int) (*entity.Attendance, error)
	GetMonthlyAttendances(employeeID int, year, month int) ([]entity.Attendance, error)
}

type attendanceService struct {
	attendanceRepo  attendance.Repository
	workHoursConfig *config.WorkHoursConfig
	timeProvider    TimeProvider
}

func NewAttendanceService(attendanceRepo attendance.Repository, workHoursConfig *config.WorkHoursConfig) AttendanceService {
	return &attendanceService{
		attendanceRepo:  attendanceRepo,
		workHoursConfig: workHoursConfig,
		timeProvider:    &RealTimeProvider{},
	}
}

// NewAttendanceServiceWithTimeProvider はテスト用にTimeProviderを注入できるコンストラクタ
func NewAttendanceServiceWithTimeProvider(attendanceRepo attendance.Repository, workHoursConfig *config.WorkHoursConfig, timeProvider TimeProvider) AttendanceService {
	return &attendanceService{
		attendanceRepo:  attendanceRepo,
		workHoursConfig: workHoursConfig,
		timeProvider:    timeProvider,
	}
}

func (s *attendanceService) ClockIn(employeeID int, latitude, longitude *float64, clockSource int) (*entity.Attendance, error) {
	today := s.timeProvider.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	// 既に打刻済みかチェック
	existing, err := s.attendanceRepo.FindByEmployeeAndDate(employeeID, todayDate)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil && existing.OpeningTime != nil {
		return nil, errors.New(errors.ErrCodeAlreadyClockedIn, "Already clocked in today")
	}

	attendanceStatus := int(attendance.DetermineClockInStatus(
		today,
		s.workHoursConfig.StartHour,
		s.workHoursConfig.StartMinute,
	))

	attendance := &entity.Attendance{
		EmployeeID:       employeeID,
		TargetDate:       todayDate,
		OpeningTime:      &today,
		AttendanceStatus: attendanceStatus,
		ClockSource:      clockSource,
		Latitude:         latitude,
		Longitude:        longitude,
	}

	if existing != nil && existing.AttendanceID > 0 {
		// 既存レコードを更新
		attendance.AttendanceID = existing.AttendanceID
		if err := s.attendanceRepo.Update(attendance); err != nil {
			return nil, err
		}
	} else {
		// 新規作成
		if err := s.attendanceRepo.Create(attendance); err != nil {
			return nil, err
		}
	}

	return attendance, nil
}

func (s *attendanceService) ClockOut(employeeID int, latitude, longitude *float64) (*entity.Attendance, error) {
	today := s.timeProvider.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	record, err := s.attendanceRepo.FindByEmployeeAndDate(employeeID, todayDate)
	if err != nil {
		return nil, errors.New(errors.ErrCodeNotClockedIn, "No clock-in record found")
	}

	if record.ClosingTime != nil {
		return nil, errors.New(errors.ErrCodeAlreadyClockedOut, "Already clocked out")
	}

	record.ClosingTime = &today

	if latitude != nil {
		record.Latitude = latitude
	}
	if longitude != nil {
		record.Longitude = longitude
	}

	// 勤務時間計算
	if record.OpeningTime != nil {
		calculator := &attendance.WorkHoursCalculator{
			BreakHours:        s.workHoursConfig.BreakHours,
			StandardWorkHours: s.workHoursConfig.StandardWorkHours,
		}

		workHours := calculator.CalculateWorkHours(*record.OpeningTime, today)
		record.WorkHours = &workHours

		overtimeHours := calculator.CalculateOvertimeHours(workHours)
		if overtimeHours > 0 {
			record.OvertimeHours = &overtimeHours
		}
	}

	// 早退判定
	if attendance.ShouldUpdateToEarlyLeave(
		today,
		attendance.AttendanceStatus(record.AttendanceStatus),
		s.workHoursConfig.EndHour,
		s.workHoursConfig.EndMinute,
	) {
		record.AttendanceStatus = int(attendance.StatusEarlyLeave)
	}

	if err := s.attendanceRepo.Update(record); err != nil {
		return nil, err
	}

	return record, nil
}

func (s *attendanceService) GetTodayAttendance(employeeID int) (*entity.Attendance, error) {
	today := s.timeProvider.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	return s.attendanceRepo.FindByEmployeeAndDate(employeeID, todayDate)
}

func (s *attendanceService) GetMonthlyAttendances(employeeID int, year, month int) ([]entity.Attendance, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)
	return s.attendanceRepo.FindByEmployeeAndDateRange(employeeID, startDate, endDate)
}
