package service

import (
	"errors"
	"time"

	"github.com/nanato-okajima/attendance_management/config"
	"github.com/nanato-okajima/attendance_management/domain"
	"github.com/nanato-okajima/attendance_management/models"
	"github.com/nanato-okajima/attendance_management/service/repository"
	"gorm.io/gorm"
)

type AttendanceService interface {
	ClockIn(employeeID int, latitude, longitude *float64, clockSource int) (*models.Attendance, error)
	ClockOut(employeeID int, latitude, longitude *float64) (*models.Attendance, error)
	GetTodayAttendance(employeeID int) (*models.Attendance, error)
	GetMonthlyAttendances(employeeID int, year, month int) ([]models.Attendance, error)
}

type attendanceService struct {
	attendanceRepo  repository.AttendanceRepository
	workHoursConfig *config.WorkHoursConfig
	timeProvider    TimeProvider
}

func NewAttendanceService(attendanceRepo repository.AttendanceRepository, workHoursConfig *config.WorkHoursConfig) AttendanceService {
	return &attendanceService{
		attendanceRepo:  attendanceRepo,
		workHoursConfig: workHoursConfig,
		timeProvider:    &RealTimeProvider{},
	}
}

// NewAttendanceServiceWithTimeProvider はテスト用にTimeProviderを注入できるコンストラクタ
func NewAttendanceServiceWithTimeProvider(attendanceRepo repository.AttendanceRepository, workHoursConfig *config.WorkHoursConfig, timeProvider TimeProvider) AttendanceService {
	return &attendanceService{
		attendanceRepo:  attendanceRepo,
		workHoursConfig: workHoursConfig,
		timeProvider:    timeProvider,
	}
}

func (s *attendanceService) ClockIn(employeeID int, latitude, longitude *float64, clockSource int) (*models.Attendance, error) {
	today := s.timeProvider.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	// 既に打刻済みかチェック
	existing, err := s.attendanceRepo.FindByEmployeeAndDate(employeeID, todayDate)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil && existing.OpeningTime != nil {
		return nil, errors.New("already clocked in today")
	}

	attendanceStatus := int(domain.DetermineClockInStatus(
		today,
		s.workHoursConfig.StartHour,
		s.workHoursConfig.StartMinute,
	))

	attendance := &models.Attendance{
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

func (s *attendanceService) ClockOut(employeeID int, latitude, longitude *float64) (*models.Attendance, error) {
	today := s.timeProvider.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)

	attendance, err := s.attendanceRepo.FindByEmployeeAndDate(employeeID, todayDate)
	if err != nil {
		return nil, errors.New("no clock-in record found")
	}

	if attendance.ClosingTime != nil {
		return nil, errors.New("already clocked out")
	}

	attendance.ClosingTime = &today

	if latitude != nil {
		attendance.Latitude = latitude
	}
	if longitude != nil {
		attendance.Longitude = longitude
	}

	// 勤務時間計算
	if attendance.OpeningTime != nil {
		calculator := domain.WorkHoursCalculator{
			BreakHours:        s.workHoursConfig.BreakHours,
			StandardWorkHours: s.workHoursConfig.StandardWorkHours,
		}

		workHours := calculator.CalculateWorkHours(*attendance.OpeningTime, today)
		attendance.WorkHours = &workHours

		overtimeHours := calculator.CalculateOvertimeHours(workHours)
		if overtimeHours > 0 {
			attendance.OvertimeHours = &overtimeHours
		}
	}

	// 早退判定
	if domain.ShouldUpdateToEarlyLeave(
		today,
		domain.AttendanceStatus(attendance.AttendanceStatus),
		s.workHoursConfig.EndHour,
		s.workHoursConfig.EndMinute,
	) {
		attendance.AttendanceStatus = int(domain.AttendanceStatusEarlyLeave)
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) GetTodayAttendance(employeeID int) (*models.Attendance, error) {
	today := s.timeProvider.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	return s.attendanceRepo.FindByEmployeeAndDate(employeeID, todayDate)
}

func (s *attendanceService) GetMonthlyAttendances(employeeID int, year, month int) ([]models.Attendance, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)
	return s.attendanceRepo.FindByEmployeeAndDateRange(employeeID, startDate, endDate)
}
