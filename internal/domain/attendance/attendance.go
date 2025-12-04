package attendance

import "time"

// AttendanceStatus は勤怠記録のステータスを表す
//
//nolint:revive // Explicit naming preferred for clarity
type AttendanceStatus int

const (
	StatusNormal     AttendanceStatus = 1 // 通常出勤
	StatusAbsent     AttendanceStatus = 2 // 欠勤
	StatusLate       AttendanceStatus = 3 // 遅刻
	StatusEarlyLeave AttendanceStatus = 4 // 早退
	StatusHoliday    AttendanceStatus = 5 // 休日出勤
)

// ClockSource は打刻元を表す
type ClockSource int

const (
	SourceApp   ClockSource = 1 // アプリ
	SourceLINE  ClockSource = 2 // LINE
	SourceAdmin ClockSource = 3 // 管理画面
)

// DetermineClockInStatus は出勤時刻に基づいて勤怠ステータスを判定する
// 出勤時刻が始業時刻より後の場合は StatusLate を返し、
// それ以外の場合は StatusNormal を返す
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
		return StatusLate
	}

	return StatusNormal
}

// ShouldUpdateToEarlyLeave はステータスを早退に更新すべきかを判定する
// 退勤時刻が終業時刻より前で、かつ現在のステータスが通常出勤の場合に true を返す
func ShouldUpdateToEarlyLeave(clockOutTime time.Time, currentStatus AttendanceStatus, workEndHour, workEndMinute int) bool {
	if currentStatus != StatusNormal {
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
