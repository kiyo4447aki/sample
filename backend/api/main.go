package main

import (
	"context"
	"log"

	"backend-proto/config"
	"backend-proto/handlers"
	"backend-proto/middleware"
	"backend-proto/utils"
	"backend-proto/utils/notifications"
	"backend-proto/utils/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err.Error())
	}

	db, err := utils.InitDB(cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("データベース接続の初期化に失敗しました: %v",err.Error())
	}

	//FCM初期化
	//TODO：クレデンシャル取得方法を変更
	fcm, err := notifications.NewFCM(ctx, cfg.GoogleCredPath)
	if err != nil {
		log.Fatalf("fcmの初期化に失敗しました：%v", err.Error())
	}

	//gcs初期化
	//FIXME：バケット名をハードコードせずconfigから取得
	gcs, err := storage.NewGCSStore(ctx, "", cfg.GoogleCredPath)
	if err != nil {
		log.Fatalf("gcsの初期化に失敗しました：%v", err.Error())
	}

	handlers.SetDB(db)
	handlers.SetConf(cfg)

	router := gin.Default()

	router.POST("/register", handlers.RegisterUserHandler)
	router.POST("/login", handlers.LoginHandler)
	router.POST("/register/device", handlers.CreateDeviceHandler)
	router.POST("/register/device/assign", handlers.AssignUserDeviceHandler)
	router.Any("/auth", handlers.AuthRequestHandler)

	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		//トークンリフレッシュ
		authGroup.POST("/refresh", handlers.TokenRefreshHandler)

		//録画映像関連
		authGroup.GET("/recordings", handlers.RecordingsHandler)
		authGroup.GET("/recordings/:recording_id", handlers.RecordingDetailHandler)
		
		//デバイス関連
		authGroup.GET("/devices/:device_id/connection-info", handlers.DeviceConnectionInfoHandler)
		authGroup.GET("/devices", handlers.UserDevicesHandler)

		//通知関連
		authGroup.POST("/notify/token/register", handlers.RegisterPushTokenHandler)

		//熱源アラート
		authGroup.POST("/alerts/:device_id/uploads", handlers.NewPostImageAlertHandler(gcs))
		authGroup.POST("/alerts/:device_id", handlers.NewPostAlertHandler(gcs, fcm))
		authGroup.GET("/alerts/:device_id", handlers.NewGetAlertHandler(gcs))

	}

	

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}