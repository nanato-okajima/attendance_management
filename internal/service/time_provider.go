package service

import "time"

// TimeProvider は現在時刻を提供するインターフェース
type TimeProvider interface {
	Now() time.Time
}

// RealTimeProvider は実際の時刻を返す実装
type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

// FixedTimeProvider はテスト用の固定時刻を返す実装
type FixedTimeProvider struct {
	FixedTime time.Time
}

func (f *FixedTimeProvider) Now() time.Time {
	return f.FixedTime
}
