package config

// DatabaseConfig はデータベース接続設定
type DatabaseConfig struct {
	User     string `envconfig:"DB_USER" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
}
