package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	LINE     LINEConfig
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

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
