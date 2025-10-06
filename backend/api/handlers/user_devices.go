package handlers

import (
	"net/http"

	"backend-proto/models"

	"github.com/gin-gonic/gin"
)

type AssignUserDeviceRequest struct {
	UserID string `json:"user_id" binding:"required"` 
	DeviceID string `json:"device_id" binding:"required"` 
}

func AssignUserDeviceHandler(c *gin.Context){
	var req AssignUserDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status": "failed",
				"error": "リクエストデータが不正です"+err.Error(),
			},
		)
		return
	}

	user, err := models.GetUserByUsername(db, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status": "failed",
				"error": "ユーザーの取得に失敗しました"+err.Error(),
			},
		)
		return
	}

	err = models.AddUserDevice(db, user.ID, req.DeviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status": "failed",
				"error": "ユーザー-デバイスの関連付けを登録できませんでした"+err.Error(),
			},
		)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
