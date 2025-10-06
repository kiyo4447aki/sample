package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DeviceConfig struct {
	RecordDir  string `yaml:"recdir"`
	Port int `yaml:"port"`
}

type Config struct {
	JanusURL string `yaml:"janus_url"`
	Devices map[string]DeviceConfig `yaml:"devices"`
	RoomID string `yaml:"room_id"`
	RoomPass string `yaml:"room_pass"`
}

func LoadConfig() (*Config, error){
	path := "./config.yaml"

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました (%s): %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗しました: %w", err)
	}

	//バリデーション
	if cfg.JanusURL == "" {
		return nil, fmt.Errorf("janus_url が設定されていません")
	}
	if len(cfg.Devices) == 0 {
		return nil, fmt.Errorf("1つ以上のdevicesを登録してください")
	}
	if cfg.RoomID == "" {
		return nil, fmt.Errorf("room_id が設定されていません")
	}
	if cfg.RoomPass == "" {
		return nil, fmt.Errorf("room_pass が設定されていません")
	}
	for id, dev := range cfg.Devices {
		if dev.RecordDir == "" {
			return nil, fmt.Errorf("devices.%s.recdir が設定されていません", id)
		}
		if dev.Port == 0 {
			return nil, fmt.Errorf("devices.%s.port が設定されていません", id)
		}
	}
	return &cfg , nil
}
