package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database  DatabaseConfig
	JWT       JWTConfig
	LINE      LINEConfig
	WorkHours WorkHoursConfig
}

type DatabaseConfig struct {
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}

type JWTConfig struct {
	Secret     string `envconfig:"JWT_SECRET" required:"true"`
	Expiration string `envconfig:"JWT_EXPIRATION" default:"24h"`
}

type LINEConfig struct {
	ChannelSecret string `envconfig:"LINE_CHANNEL_SECRET"`
	ChannelToken  string `envconfig:"LINE_CHANNEL_TOKEN"`
}

type WorkHoursConfig struct {
	StartHour         int     `envconfig:"WORK_START_HOUR" default:"9"`
	StartMinute       int     `envconfig:"WORK_START_MINUTE" default:"0"`
	EndHour           int     `envconfig:"WORK_END_HOUR" default:"18"`
	EndMinute         int     `envconfig:"WORK_END_MINUTE" default:"0"`
	BreakHours        float64 `envconfig:"BREAK_HOURS" default:"1"`
	StandardWorkHours float64 `envconfig:"STANDARD_WORK_HOURS" default:"8"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
