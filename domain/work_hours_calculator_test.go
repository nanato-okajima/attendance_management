package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateWorkHours(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		openingTime time.Time
		closingTime time.Time
		expected    float64
		description string
	}{
		{
			name:        "休憩時間あり",
			openingTime: time.Date(2025, 11, 30, 9, 0, 0, 0, time.Local),
			closingTime: time.Date(2025, 11, 30, 18, 0, 0, 0, time.Local),
			expected:    8.0,
			description: "9時間勤務 - 1時間休憩 = 8時間",
		},
		{
			name:        "休憩時間なし（短時間勤務）",
			openingTime: time.Date(2025, 11, 30, 9, 0, 0, 0, time.Local),
			closingTime: time.Date(2025, 11, 30, 9, 30, 0, 0, time.Local),
			expected:    0.5,
			description: "0.5時間勤務（休憩時間より短いため休憩は引かれない）",
		},
		{
			name:        "休憩時間ちょうど",
			openingTime: time.Date(2025, 11, 30, 9, 0, 0, 0, time.Local),
			closingTime: time.Date(2025, 11, 30, 10, 0, 0, 0, time.Local),
			expected:    1.0,
			description: "1時間勤務（休憩時間と同じなので休憩は引かれない）",
		},
		{
			name:        "長時間勤務",
			openingTime: time.Date(2025, 11, 30, 9, 0, 0, 0, time.Local),
			closingTime: time.Date(2025, 11, 30, 20, 0, 0, 0, time.Local),
			expected:    10.0,
			description: "11時間勤務 - 1時間休憩 = 10時間",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calculator := WorkHoursCalculator{
				BreakHours:        1.0,
				StandardWorkHours: 8.0,
			}

			workHours := calculator.CalculateWorkHours(tt.openingTime, tt.closingTime)
			assert.Equal(t, tt.expected, workHours, tt.description)
		})
	}
}

func TestCalculateOvertimeHours(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		workHours   float64
		expected    float64
		description string
	}{
		{
			name:        "残業あり",
			workHours:   10.0,
			expected:    2.0,
			description: "10時間 - 8時間（標準） = 2時間残業",
		},
		{
			name:        "残業なし",
			workHours:   7.0,
			expected:    0.0,
			description: "標準勤務時間以下なので残業なし",
		},
		{
			name:        "標準勤務時間ちょうど",
			workHours:   8.0,
			expected:    0.0,
			description: "ちょうど標準勤務時間なので残業なし",
		},
		{
			name:        "わずかな残業",
			workHours:   8.5,
			expected:    0.5,
			description: "8.5時間 - 8時間（標準） = 0.5時間残業",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calculator := WorkHoursCalculator{
				BreakHours:        1.0,
				StandardWorkHours: 8.0,
			}

			overtimeHours := calculator.CalculateOvertimeHours(tt.workHours)
			assert.Equal(t, tt.expected, overtimeHours, tt.description)
		})
	}
}
