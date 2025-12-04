package config

import "fmt"

// Validate は設定値のバリデーションを行う
func (c *Config) Validate() error {
	// データベース設定のバリデーション
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	// JWT設定のバリデーション
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// 勤務時間設定のバリデーション
	if c.WorkHours.StandardWorkHours <= 0 {
		return fmt.Errorf("standard work hours must be positive")
	}
	if c.WorkHours.BreakHours < 0 {
		return fmt.Errorf("break hours cannot be negative")
	}

	// サーバー設定のバリデーション
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	return nil
}
