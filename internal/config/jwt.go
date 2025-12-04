package config

// JWTConfig はJWT認証設定
type JWTConfig struct {
	Secret     string `envconfig:"JWT_SECRET" required:"true"`
	Expiration string `envconfig:"JWT_EXPIRATION" default:"24h"`
}
