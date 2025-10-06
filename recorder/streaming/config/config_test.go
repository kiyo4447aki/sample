// config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func setupTestConfig(t *testing.T, yamlContent string) func() {
	tmp := t.TempDir()
	workDir := filepath.Join(tmp, "work")
	require.NoError(t, os.Mkdir(workDir, 0755))

	configPath := filepath.Join(workDir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(yamlContent), 0644))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workDir))

	return func() {
		_ = os.Chdir(origDir)
	}
}

func TestLoadConfig_Success(t *testing.T) {
	yaml := `
janus_url: "wss://example.com"
devices:
  dev1:
    recdir: "/tmp/rec1"
    port: 47000
room_id: "roomX"
room_pass: "passX"
`
	cleanup := setupTestConfig(t, yaml)
	defer cleanup()

	cfg, err := LoadConfig()
	require.NoError(t, err)

	assert.Equal(t, "wss://example.com", cfg.JanusURL)
	assert.Equal(t, "roomX", cfg.RoomID)
	assert.Equal(t, "passX", cfg.RoomPass)

	require.Len(t, cfg.Devices, 1)
	dev, ok := cfg.Devices["dev1"]
	require.True(t, ok)
	assert.Equal(t, "/tmp/rec1", dev.RecordDir)
	assert.Equal(t, 47000, dev.Port)
}

func TestLoadConfig_NoFile(t *testing.T) {
	tmp := t.TempDir()
	workDir := filepath.Join(tmp, "work")
	require.NoError(t, os.Mkdir(workDir, 0755))

	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workDir))
	defer func() { _ = os.Chdir(origDir) }()

	_, err = LoadConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "設定ファイルの読み込みに失敗しました")
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	invalid := `: not valid yaml`
	cleanup := setupTestConfig(t, invalid)
	defer cleanup()

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "設定ファイルのパースに失敗しました")
}

func TestLoadConfig_NoJanusURL(t *testing.T) {
	yaml := `
janus_url: ""
devices:
  d:
    recdir: "/tmp"
room_id: "r"
room_pass: "p"
`
	cleanup := setupTestConfig(t, yaml)
	defer cleanup()

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Equal(t, "janus_url が設定されていません", err.Error())
}

func TestLoadConfig_NoDevices(t *testing.T) {
	yaml := `
janus_url: "wss://x"
devices: {}
room_id: "r"
room_pass: "p"
`
	cleanup := setupTestConfig(t, yaml)
	defer cleanup()

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Equal(t, "1つ以上のdevicesを登録してください", err.Error())
}

func TestLoadConfig_NoRoomID(t *testing.T) {
	yaml := `
janus_url: "wss://x"
devices:
  d:
    recdir: "/tmp"
room_id: ""
room_pass: "p"
`
	cleanup := setupTestConfig(t, yaml)
	defer cleanup()

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Equal(t, "room_id が設定されていません", err.Error())
}

func TestLoadConfig_NoRoomPass(t *testing.T) {
	yaml := `
janus_url: "wss://x"
devices:
  d:
    recdir: "/tmp"
room_id: "r"
room_pass: ""
`
	cleanup := setupTestConfig(t, yaml)
	defer cleanup()

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Equal(t, "room_pass が設定されていません", err.Error())
}

func TestLoadConfig_DeviceMissingRecDir(t *testing.T) {
	yaml := `
janus_url: "wss://x"
devices:
  d1:
    recdir: ""
room_id: "r"
room_pass: "p"
`
	cleanup := setupTestConfig(t, yaml)
	defer cleanup()

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Equal(t, "devices.d1.recdir が設定されていません", err.Error())
}

