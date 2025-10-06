package main

import (
	"api/config"
	"api/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	router := gin.Default()
	router.GET("/records", handlers.ListRecordsHandler(cfg))

	router.Run(cfg.ListenAddress)
}







