package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//正常系
func TestLoadConfig_SuccessWithCustomListenAddress(t *testing.T) {
	t.Setenv("BASE_PATH", "/var/data/records")
	t.Setenv("BASE_URL", "http://localhost:8000/play")
	t.Setenv("LISTEN_ADDRESS", "0.0.0.0:8080")

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "/var/data/records", cfg.BasePath)
	assert.Equal(t, "http://localhost:8000/play", cfg.BaseUrl)
	assert.Equal(t, "0.0.0.0:8080", cfg.ListenAddress)
}

//正常系：デフォルトのアドレス使用時
func TestLoadConfig_SuccessDefaultListenAddress(t *testing.T) {
	t.Setenv("BASE_PATH", "/tmp/rec")
	t.Setenv("BASE_URL", "https://example.com/stream")
	t.Setenv("LISTEN_ADDRESS", "")

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "/tmp/rec", cfg.BasePath)
	assert.Equal(t, "https://example.com/stream", cfg.BaseUrl)
	assert.Equal(t, "127.0.0.1:5000", cfg.ListenAddress)
}

//異常系：BASE_PATHが未設定のとき適切なエラーが返る
func TestLoadConfig_MissingBasePath(t *testing.T) {
	t.Setenv("BASE_PATH", "")
	t.Setenv("BASE_URL", "http://localhost/play")

	cfg, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "環境変数 BASE_PATH がセットされていません")
}

//異常系：BASE_URLが未設定のとき適切なエラーが返る
func TestLoadConfig_MissingBaseUrl(t *testing.T) {
	t.Setenv("BASE_PATH", "/records")
	t.Setenv("BASE_URL", "")

	cfg, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "環境変数 BASE_ URL がセットされていません")
}
