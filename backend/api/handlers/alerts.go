package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend-proto/models"
	"backend-proto/utils/notifications"
	"backend-proto/utils/storage"

	gcs "cloud.google.com/go/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateUploadReq struct {
	MimeType     string     `json:"mime_type" binding:"required"`
	Width        *int       `json:"width,omitempty"`
	Height       *int       `json:"height,omitempty"`
	SHA256       *string    `json:"sha256,omitempty"`
	OccurredAt   string     `json:"occurred_at" binding:"required"`
	EventID      string     `json:"event_id" binding:"required"`
}

type CreateUploadResp struct {
	Status    string `json:"status"`
	ObjectKey string `json:"object_key"`
	PutURL    string `json:"put_url"`
	ExpiresIn int    `json:"expires_in"`
}

// POST /alerts/:device_id/uploads
func NewPostImageAlertHandler(store *storage.GCSStore)gin.HandlerFunc{
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		isAuth := authWithErrorFromCtx(c, deviceID)
		if !isAuth {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error", 
				"error":"権限がありません",
			})
			return 
		}

		var req CreateUploadReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error" : "リクエストの形式が不正です",
			})
			return 
		}

		occurredAt, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error": "occurred_at の形式が不正です",
			})
			return 
		}
		occurredAt = occurredAt.UTC()

		switch req.MimeType {
        case "image/jpeg", "image/png":
        default:
            c.JSON(http.StatusBadRequest, gin.H{"status": "error","error":"サポートされていないmime_typeです"})
            return
        }

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
        defer cancel()
		objectKey := getObjectKey(deviceID, req.EventID, occurredAt, req.MimeType)
		var existing models.AlertMedia
		err = db.WithContext(ctx).Where(
			"object_key = ? AND committed = TRUE", objectKey,
		).First(&existing).Error
		if err == nil {
			c.JSON(
				http.StatusConflict, gin.H{
					"status": "error",
					"error": "すでにコミットされたアラートです",
			}) 
			return 
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "メディア情報の確認に失敗しました",
			})
			return
		}
		ttl := 5*time.Minute
		putURL, err := store.PresignedPutURL(objectKey, req.MimeType, ttl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "アップロード用URLを取得できませんでした",
			})
			return 
		}

		media := models.AlertMedia{
			ObjectKey: objectKey,
			MimeType: req.MimeType,
			Width: req.Width,
			Height: req.Height,
			SHA256: req.SHA256,
			Committed: false,
		}

		if err := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&media).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "メディアの情報を取得できませんでした",
			})
			return 
		}

		c.JSON(http.StatusOK, CreateUploadResp{
			Status: "success",
			ObjectKey: objectKey,
			PutURL: putURL,
			ExpiresIn: int(ttl/time.Second),
		})
	}
}


// POST /alerts/:device_id
type ThermalAlertReq struct {
	Severity   string   `json:"severity"`
	TempC      *float64 `json:"temp_c"`
	OccurredAt string   `json:"occurred_at" binding:"required"`
	EventID    string   `json:"event_id" binding:"required"`
	Tags       []string `json:"tags"`
	MimeType   string  `json:"mime_type" binding:"required"`
	ObjectKey  string  `json:"object_key" binding:"required"`
}

type ThermalAlertResp struct {
	Status    string `json:"status"`
	AlertID uint64 `json:"alert_id"`
}

func NewPostAlertHandler(store *storage.GCSStore, fcm *notifications.FCM)gin.HandlerFunc{
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")
		isAuth := authWithErrorFromCtx(c, deviceID)
		if !isAuth {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error", 
				"error":"権限がありません",
			})
			return 
		}

		var req ThermalAlertReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error" : "リクエストの形式が不正です",
			})
			return
		}

		occurredAt, err := time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error": "occurred_atの形式が不正です",
			})
			return 
		}
		occurredAt = occurredAt.UTC()

		var tempCValue float64
		if req.TempC != nil {
			tempCValue = *req.TempC
		} else {
			tempCValue = 0.0
		}

		alert := models.Alert{
			DeviceID: deviceID,
			EventID: req.EventID,
			Severity: ifEmpty(req.Severity, "info"),
			TempC: req.TempC,
			OccurredAt: occurredAt,
		}

		var objectKey string
		if req.ObjectKey != ""{
			if !validateObjectKey(req.ObjectKey, deviceID, req.EventID, req.MimeType){
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "error",
					"error": "invalid objectkey",
				})
				return 
			}
			objectKey = req.ObjectKey
		} else {
			objectKey = getObjectKey(deviceID, req.EventID, occurredAt, req.MimeType)
		}

		ctxTimeout, cancelTimeout := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancelTimeout()

		attrs, err := store.Stat(ctxTimeout, objectKey)
		notfound := errors.Is(err, gcs.ErrObjectNotExist)
		switch {
		case err == nil:
		case notfound:
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"status":"error",
				"error":"オブジェクト情報の取得中にタイムアウトしました",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "オブジェクト情報を取得できませんでした",
			})
			return 
		}
		
		

		conflictErr := fmt.Errorf("event_id conflicts")

		txErr := db.WithContext(ctxTimeout).Transaction(func(tx *gorm.DB) error {
			res := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "event_id"}},
				DoNothing: true,
			}).Create(&alert)

				if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected == 0 {
				return conflictErr
			}

			
			upd := map[string]any{"alert_id": alert.ID}
			if !notfound{
				upd["committed"] = true
				upd["size_bytes"] = attrs.Size
			}

			q := tx.Model(&models.AlertMedia{}).Where(
				"object_key = ?", objectKey,
			).Updates(upd)
			if q.Error != nil {
				return q.Error
			}

			if q.RowsAffected == 0 {
				media := models.AlertMedia{
					AlertID:   &alert.ID,
					ObjectKey: objectKey,
					MimeType:  req.MimeType,
					Committed: !notfound,
				}
				if !notfound {
					sz := attrs.Size
					media.SizeBytes = &sz
				}
				if err := tx.Clauses(clause.OnConflict{
					DoNothing: true,
					}).Create(&media).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if txErr != nil{
			if errors.Is(txErr, conflictErr) {
				c.JSON(http.StatusConflict, gin.H{
					"status": "error",
					"error":"event_idが重複しています",
				})
				return 
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":"アラートの登録に失敗しました",
			})
			return
		}

		
		tokens, err := models.GetAvailableTokensByDeviceID(db.WithContext(ctxTimeout), deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "プッシュトークンを取得できませんでした",
			})
			return 
		}

		
		payload := notifications.ThermalAlertPayload{
			DeviceID: deviceID,
			Severity: alert.Severity,
			TempC: tempCValue,
			OccurredAt: occurredAt.Format(time.RFC3339),
			AlertID: strconv.FormatUint(alert.ID, 10),
		}

		ctxCancel, cancel := context.WithCancel(c.Request.Context())
		defer cancel()
		err = fcm.SendThermalAlert(ctxCancel, db, payload, tokens)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "アラートの配信に失敗しました",
			})
			return 
		}
		c.JSON(http.StatusOK, ThermalAlertResp{
			Status: "success",
			AlertID: alert.ID,
		})
	}
}

// GET /alerts/:id
type GetAlertResp struct {
	ID         uint64         `json:"id"`
	DeviceID   string         `json:"device_id"`
	Severity   string         `json:"severity"`
	OccurredAt time.Time      `json:"occurred_at"`
	Media      MediaInfoResp  `json:"media,omitempty"`
}

type MediaInfoResp struct {
	IsExists  bool   `json:"is_exists"`
	MimeType  *string `json:"mime_type"`
	SizeBytes *int64 `json:"size_bytes"`
	Width     *int   `json:"width"`
	Height    *int   `json:"height"`
	SHA256    *string `json:"sha256"`
	URL       *string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}

func NewGetAlertHandler(store *storage.GCSStore) gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		aid, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error": "not found",
			})
			return
		}

		var alert models.Alert
		if err := db.WithContext(ctx).First(&alert, aid).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound){
				c.JSON(http.StatusNotFound, gin.H{
					"status": "error",
					"error": "not found",
				})
				return 
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error": "アラート情報を取得できませんでした",
			})
			return 
		}

		isAuth := authWithErrorFromCtx(c, alert.DeviceID)
		if !isAuth {
			c.JSON(http.StatusForbidden, gin.H{
				"status": "error", 
				"error":"権限がありません",
			})
			return 
		}

		var media models.AlertMedia
		if err := db.WithContext(ctx).Where(
			"alert_id = ? AND committed = TRUE", alert.ID,
		).First(&media).Error; err != nil {
			//画像がない場合のアラート詳細
			c.JSON(http.StatusOK, GetAlertResp{
				ID: alert.ID,
				DeviceID: alert.DeviceID,
				Severity: alert.Severity,
				OccurredAt: alert.OccurredAt,
				Media: MediaInfoResp{
					IsExists: false,
				},
			})
			return 
		}

		ttl := 3*time.Minute

		url, err := store.PresignedGetURL(media.ObjectKey, ttl)
		if err != nil {
			c.JSON(http.StatusOK, GetAlertResp{
				ID: alert.ID,
				DeviceID: alert.DeviceID,
				Severity: alert.Severity,
				OccurredAt: alert.OccurredAt,
				Media: MediaInfoResp{
					IsExists: true,
				},
			})
			return 
		}

		c.JSON(http.StatusOK, GetAlertResp{
			ID: alert.ID,
			DeviceID: alert.DeviceID,
			Severity: alert.Severity,
			OccurredAt: alert.OccurredAt,
			Media: MediaInfoResp{
				IsExists: true,
				MimeType: &media.MimeType,
				SizeBytes: media.SizeBytes,
				Width: media.Width,
				Height: media.Height,
				SHA256: media.SHA256,
				URL: &url,
				ExpiresIn: int(ttl/time.Second),
			},
		})
	}
}