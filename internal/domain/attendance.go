package domain

import "time"

// AttendanceStatus は勤怠記録のステータスを表す
type AttendanceStatus int

const (
	AttendanceStatusNormal     AttendanceStatus = 1 // 通常出勤
	AttendanceStatusAbsent     AttendanceStatus = 2 // 欠勤
	AttendanceStatusLate       AttendanceStatus = 3 // 遅刻
	AttendanceStatusEarlyLeave AttendanceStatus = 4 // 早退
	AttendanceStatusHoliday    AttendanceStatus = 5 // 休日出勤
)

// ClockSource は打刻元を表す
type ClockSource int

const (
	ClockSourceApp   ClockSource = 1 // アプリ
	ClockSourceLINE  ClockSource = 2 // LINE
	ClockSourceAdmin ClockSource = 3 // 管理画面
)

// DetermineClockInStatus は出勤時刻に基づいて勤怠ステータスを判定する
// 出勤時刻が始業時刻より後の場合は AttendanceStatusLate を返し、
// それ以外の場合は AttendanceStatusNormal を返す
func DetermineClockInStatus(clockInTime time.Time, workStartHour, workStartMinute int) AttendanceStatus {
	standardTime := time.Date(
		clockInTime.Year(),
		clockInTime.Month(),
		clockInTime.Day(),
		workStartHour,
		workStartMinute,
		0, 0,
		time.Local,
	)

	if clockInTime.After(standardTime) {
		return AttendanceStatusLate
	}

	return AttendanceStatusNormal
}

// ShouldUpdateToEarlyLeave はステータスを早退に更新すべきかを判定する
// 退勤時刻が終業時刻より前で、かつ現在のステータスが通常出勤の場合に true を返す
func ShouldUpdateToEarlyLeave(clockOutTime time.Time, currentStatus AttendanceStatus, workEndHour, workEndMinute int) bool {
	if currentStatus != AttendanceStatusNormal {
		return false
	}

	standardEndTime := time.Date(
		clockOutTime.Year(),
		clockOutTime.Month(),
		clockOutTime.Day(),
		workEndHour,
		workEndMinute,
		0, 0,
		time.Local,
	)

	return clockOutTime.Before(standardEndTime)
}

// WorkHoursCalculator は勤務時間計算のロジックを持つ
type WorkHoursCalculator struct {
	BreakHours        float64
	StandardWorkHours float64
}

// CalculateWorkHours は勤務時間を計算する
// 出勤時刻と退勤時刻から勤務時間を算出し、休憩時間を差し引く
func (c *WorkHoursCalculator) CalculateWorkHours(openingTime, closingTime time.Time) float64 {
	duration := closingTime.Sub(openingTime)
	workHours := duration.Hours()

	// 休憩時間を引く
	if workHours > c.BreakHours {
		workHours -= c.BreakHours
	}

	return workHours
}

// CalculateOvertimeHours は残業時間を計算する
// 勤務時間が標準勤務時間を超える場合、その差分を残業時間として返す
func (c *WorkHoursCalculator) CalculateOvertimeHours(workHours float64) float64 {
	if workHours > c.StandardWorkHours {
		return workHours - c.StandardWorkHours
	}
	return 0
}
