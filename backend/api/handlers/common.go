package handlers

import (
	"net/http"

	"backend-proto/config"
	"backend-proto/models"
	"backend-proto/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//TODO dbをパッケージ内の変数からDIに変更
var db *gorm.DB

func SetDB(database *gorm.DB){
	db = database
}

var conf *config.Config

func SetConf(cfg *config.Config){
	conf = cfg
}

//ユーザーがデバイス対してアクセス権を持つか確認する関数
func hasDeviceAccess(userID int, deviceID string)bool{
	userDevices, err := models.GetUserDevices(db, uint(userID))
	if err != nil {
		return false
	}

	for _, ud := range userDevices{
		if ud.DeviceID == deviceID{
			return true
		}
	}
	return false
}



//ginコンテキストからデバイスへのアクセス権をチェック
//認証不可の場合エラーをレスポンス
func authWithErrorFromCtx(c *gin.Context, deviceID string)bool{
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error": "デバイスIDが指定されていません",
		})
		return false
	}

	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error", 
			"error": "認証情報が見つかりません",
		})
		return false
	}

	claims, ok := claimsInterface.(*utils.JwtClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error", 
			"error": "認証情報の形式が正しくありません",
		})
		return false
	}

	if !hasDeviceAccess(claims.UserID, deviceID) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error": "指定されたデバイスにアクセスする権限がありません",
		})
		return false
	}
	return true
}