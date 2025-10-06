package handlers

import (
	"net/http"
	"time"

	"backend-proto/models"
	"backend-proto/utils"

	"github.com/gin-gonic/gin"
)

type RegisterDeviceRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	JanusPassword string `json:"janus_password" binding:"required"`
	DeviceName string `json:"device_name" binding:"required"`
	Location string `json:"location"`
	Status string `json:"status"`
}

func DeviceConnectionInfoHandler(c *gin.Context){
	//デバイスIDを取得
	deviceID := c.Param("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": "デバイスIDが指定されていません",
		})
		return
	}

	//JWTからclaimsを取得
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error", 
			"error": "認証情報が見つかりません",
		})
		return
	}

	claims, ok := claimsInterface.(*utils.JwtClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error", 
			"error": "認証情報の形式が正しくありません",
		})
		return
	}

	//デバイスへのアクセス権を確認
	if !hasDeviceAccess(claims.UserID, deviceID) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error": "指定されたデバイスにアクセスする権限がありません",
		})
		return
	}

	//デバイス情報の取得
	device, err := models.GetDeviceByID(db, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": "デバイス情報の取得に失敗しました" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"device_id": device.DeviceID,
		"janus_ws": conf.JanusWS,
		"janus_password": device.JanusPassword, 
	})


}

func UserDevicesHandler(c *gin.Context){
	//JWTからclaims、ユーザーIDを取得
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error", 
			"error": "認証情報が見つかりません",
		})
		return
	}

	claims, ok := claimsInterface.(*utils.JwtClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error", 
			"error": "認証情報の形式が正しくありません",
		})
		return
	}

	userID := claims.UserID

	//ユーザーに紐づいているデバイスを取得
	userDevices, err := models.GetUserDevices(db, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": "デバイス情報の取得に失敗しました" + err.Error(),
		})
		return
	}

	var deviceIDs []string
	for _, ud := range userDevices{
		deviceIDs = append(deviceIDs, ud.DeviceID)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"devices": deviceIDs,
	})

}

func CreateDeviceHandler(c *gin.Context){
	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status": "failed",
				"error": "リクエストデータが不正です"+err.Error(),
			},
		)
		return
	}

	now := time.Now().UTC()
	device := models.Device{
		DeviceID: req.DeviceID,
		JanusPassword: req.JanusPassword,
		DeviceName: req.DeviceName,
		Location: req.Location,
		Status: req.Status,
		LastCommunication: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := models.CreateDevice(db, &device); err != nil {
		c.JSON(http.StatusInternalServerError, 
			gin.H{
				"status": "failed",
				"error": "デバイスの登録に失敗しました"+err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"device": device,
	})
}