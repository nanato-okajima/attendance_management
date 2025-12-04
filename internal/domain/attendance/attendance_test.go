package attendance

import (
	"testing"
	"time"
)

func TestDetermineClockInStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		clockInTime time.Time
		startHour   int
		startMinute int
		expected    AttendanceStatus
	}{
		{
			name:        "定刻前（8:59）",
			clockInTime: time.Date(2025, 11, 27, 8, 59, 0, 0, time.Local),
			startHour:   9,
			startMinute: 0,
			expected:    StatusNormal,
		},
		{
			name:        "定刻ちょうど（9:00）",
			clockInTime: time.Date(2025, 11, 27, 9, 0, 0, 0, time.Local),
			startHour:   9,
			startMinute: 0,
			expected:    StatusNormal,
		},
		{
			name:        "遅刻（9:01）",
			clockInTime: time.Date(2025, 11, 27, 9, 1, 0, 0, time.Local),
			startHour:   9,
			startMinute: 0,
			expected:    StatusLate,
		},
		{
			name:        "分指定あり_定刻前（9:29）",
			clockInTime: time.Date(2025, 11, 27, 9, 29, 0, 0, time.Local),
			startHour:   9,
			startMinute: 30,
			expected:    StatusNormal,
		},
		{
			name:        "分指定あり_遅刻（9:31）",
			clockInTime: time.Date(2025, 11, 27, 9, 31, 0, 0, time.Local),
			startHour:   9,
			startMinute: 30,
			expected:    StatusLate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			status := DetermineClockInStatus(tt.clockInTime, tt.startHour, tt.startMinute)
			if status != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, status)
			}
		})
	}
}

func TestShouldUpdateToEarlyLeave(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		clockOutTime   time.Time
		currentStatus  AttendanceStatus
		endHour        int
		endMinute      int
		expectedResult bool
	}{
		{
			name:           "通常_定刻前（17:59）",
			clockOutTime:   time.Date(2025, 11, 27, 17, 59, 0, 0, time.Local),
			currentStatus:  StatusNormal,
			endHour:        18,
			endMinute:      0,
			expectedResult: true,
		},
		{
			name:           "通常_定刻後（18:01）",
			clockOutTime:   time.Date(2025, 11, 27, 18, 1, 0, 0, time.Local),
			currentStatus:  StatusNormal,
			endHour:        18,
			endMinute:      0,
			expectedResult: false,
		},
		{
			name:           "通常_定刻ちょうど（18:00）",
			clockOutTime:   time.Date(2025, 11, 27, 18, 0, 0, 0, time.Local),
			currentStatus:  StatusNormal,
			endHour:        18,
			endMinute:      0,
			expectedResult: false,
		},
		{
			name:           "遅刻_定刻前（17:59）",
			clockOutTime:   time.Date(2025, 11, 27, 17, 59, 0, 0, time.Local),
			currentStatus:  StatusLate,
			endHour:        18,
			endMinute:      0,
			expectedResult: false,
		},
		{
			name:           "分指定あり_定刻前（18:29）",
			clockOutTime:   time.Date(2025, 11, 27, 18, 29, 0, 0, time.Local),
			currentStatus:  StatusNormal,
			endHour:        18,
			endMinute:      30,
			expectedResult: true,
		},
		{
			name:           "分指定あり_定刻後（18:31）",
			clockOutTime:   time.Date(2025, 11, 27, 18, 31, 0, 0, time.Local),
			currentStatus:  StatusNormal,
			endHour:        18,
			endMinute:      30,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ShouldUpdateToEarlyLeave(tt.clockOutTime, tt.currentStatus, tt.endHour, tt.endMinute)
			if result != tt.expectedResult {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}
