package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresDSN string // PostgreSQLのDSN
	JWTSecret   string // JWT作成時の秘密鍵
	JanusWS   string // JanusのWebSocketサーバのURL
	GoogleCredPath string
	Port string

}

func LoadConfig() (*Config, error) {
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		return nil, fmt.Errorf("環境変数 POSTGRES_DSN がセットされていません")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("環境変数 JWT_SECRET がセットされていません")
	}
	janusWS := os.Getenv("JANUS_WS")
	if janusWS == "" {
		return nil, fmt.Errorf("環境変数 JANUS_WS がセットされていません")
	}
	googleCredPath := os.Getenv("GOOGLE_CRED_PATH")
	if googleCredPath == "" {
		return nil, fmt.Errorf("環境変数 GOOGLE_CRED_PATH がセットされていません")
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		PostgresDSN: postgresDSN,
		JWTSecret:   jwtSecret,
		JanusWS:   janusWS,
		GoogleCredPath: googleCredPath,
		Port: port,
	}, nil
}