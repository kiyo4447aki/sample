package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Host           string
	Port           int
	Width          int
	Height         int
	Fps            int
	Bitrate        int
	NormalDev      string
	NightDev       string
	Tz             string
	NightVisionTime  string
	SwitchEvery    time.Duration
} 

func getenvInt(key string) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	log.Fatalf("環境変数%vの取得に失敗しました", key)
	return 0
}

func getenv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("環境変数%vの取得に失敗しました", key)
	return ""
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
	log.Fatalf(".envを読み込めませんでした：%v", err)
	}
	return &Config{
		Host:            getenv("UDP_HOST"),
		Port:            getenvInt("UDP_PORT"),
		Width:           getenvInt("WIDTH"),
		Height:          getenvInt("HEIGHT"),
		Fps:             getenvInt("FPS"),
		Bitrate:         getenvInt("BITRATE"),
		NormalDev:       getenv("NORMAL_CAMERA_PATH"),  //通常カメラのデバイスファイルパス
		NightDev:        getenv("NIGHT_CAMERA_PATH"),   //暗視カメラのデバイスファイルパス
		Tz:              getenv("TZ"),                  //Asia/Tokyo
		NightVisionTime: getenv("NIGHT_VISION_TIME"),   //暗視カメラの時間帯 hh:mm-hh:mm
		SwitchEvery:     time.Duration(getenvInt("SWITCH_EVERY_SEC")),   //カメラ切り替えをチェックするインターバル
	}
}
