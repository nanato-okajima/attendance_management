package config

// ServerConfig はサーバー設定
type ServerConfig struct {
	Port            string `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout     int    `envconfig:"SERVER_READ_TIMEOUT" default:"10"`
	WriteTimeout    int    `envconfig:"SERVER_WRITE_TIMEOUT" default:"10"`
	ShutdownTimeout int    `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"5"`
}
