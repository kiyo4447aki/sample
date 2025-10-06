package handlers

import (
	"net/http"
	"time"

	"backend-proto/utils"

	"github.com/gin-gonic/gin"
)


func RecordingsHandler(c *gin.Context){
	//クエリパラメータからデバイスIDを取得
	deviceID := c.Query("device_id")
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

	//レコーダーに問い合わせ、録画一覧を取得
	recordings, err := queryRecordingsFromDevice(deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": "録画映像一覧の取得に失敗しました: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"recordings": recordings,
	})

}

func RecordingDetailHandler(c *gin.Context){
	// パスパラメータから録画IDを取得
	recordingID := c.Param("recording_id")
	//クエリパラメータからデバイスIDを取得
	deviceID := c.Query("device_id")
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

	//レコーダーに問い合わせ、録画の詳細情報を取得
	detail, err := queryRecordingDetailFromDevice(deviceID, recordingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error": "録画映像の詳細の取得に失敗しました: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"recording": detail,
	})

}

//TODO:録画ファイル取得機能の実装
func queryRecordingsFromDevice(deviceID string) ([]map[string]interface{}, error) {
	//ダミーデータ
	recordings := []map[string]interface{}{
		{
			"id":          "rec1",
			"filename":    "video1.mp4",
			"recorded_at": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"duration":    3600,
			"device_id":   deviceID,
		},
		{
			"id":          "rec2",
			"filename":    "video2.mp4",
			"recorded_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"duration":    1800,
			"device_id":   deviceID,
		},
	}
	return recordings, nil
}

func queryRecordingDetailFromDevice(deviceID, recordingID string) (map[string]interface{}, error) {
	//ダミーデータ
	detail := map[string]interface{}{
		"id":          recordingID,
		"filename":    "video1.mp4",
		"recorded_at": time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"duration":    3600,
		"file_path":   "/path/to/video1.mp4",
		"device_id":   deviceID,
	}
	return detail, nil
}
