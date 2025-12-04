package attendance

import "time"

// WorkHoursCalculator は勤務時間計算のロジックを持つドメインサービス
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
