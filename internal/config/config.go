package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config はアプリケーション全体の設定
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	LINE      LINEConfig
	WorkHours WorkHoursConfig
}

// LINEConfig はLINE連携設定
type LINEConfig struct {
	ChannelSecret string `envconfig:"LINE_CHANNEL_SECRET"`
	ChannelToken  string `envconfig:"LINE_CHANNEL_TOKEN"`
}

// WorkHoursConfig は勤務時間設定
type WorkHoursConfig struct {
	StartHour         int     `envconfig:"WORK_START_HOUR" default:"9"`
	StartMinute       int     `envconfig:"WORK_START_MINUTE" default:"0"`
	EndHour           int     `envconfig:"WORK_END_HOUR" default:"18"`
	EndMinute         int     `envconfig:"WORK_END_MINUTE" default:"0"`
	BreakHours        float64 `envconfig:"BREAK_HOURS" default:"1"`
	StandardWorkHours float64 `envconfig:"STANDARD_WORK_HOURS" default:"8"`
}

// Load は環境変数から設定を読み込む
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	// バリデーション
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
