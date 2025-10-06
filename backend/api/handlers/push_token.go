package handlers

import (
	"net/http"
	"time"

	"backend-proto/models"
	"backend-proto/utils"

	"github.com/gin-gonic/gin"
)

type RegisterPushTokenReq struct {
	Token       string `json:"token" binding:"required"`
	Platform    string `json:"platform" binding:"required,oneof=ios android"`
}

func RegisterPushTokenHandler(c *gin.Context){
	var req RegisterPushTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status": "failed",
				"error": "リクエストデータが不正です"+err.Error(),
			})
		return
	}

	//claimsを取得
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error": "認証情報が見つかりません",
		})
		return
	}

	claims, ok := claimsInterface.(*utils.JwtClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error": "認証情報の形式が正しくありません",
		})
		return
	}

	userId := claims.UserID

	rec := models.PushToken{
		UserID: uint(userId),
		Token: req.Token,
		Platform: req.Platform,
		LastSeenAt: time.Now(),
		Revoked: false,
	}

	if err := models.RegisterToken(db, rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error": "プッシュトークンの登録に失敗しました",
		})
		return
	}

	c.Status(http.StatusNoContent)
}