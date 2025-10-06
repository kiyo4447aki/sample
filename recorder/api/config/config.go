package config

import (
	"fmt"
	"os"
)

type Config struct {
	BasePath string          //録画ファイル保存ディレクトリ
	BaseUrl string           //録画視聴URL
	ListenAddress string     //バインドするアドレス
}

func LoadConfig() (*Config, error) {
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		return nil, fmt.Errorf("環境変数 BASE_PATH がセットされていません")
	}

	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return nil, fmt.Errorf("環境変数 BASE_ URL がセットされていません")
	}

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = "127.0.0.1:5000"
	}

	return &Config{
		BasePath: basePath,
		BaseUrl: baseUrl,
		ListenAddress: listenAddress,
	}, nil
}
