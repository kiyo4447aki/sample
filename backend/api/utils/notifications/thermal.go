package notifications

import (
	"context"
	"errors"
	"fmt"
	"time"

	"firebase.google.com/go/v4/messaging"
	"gorm.io/gorm"
)

type ThermalAlertPayload struct {
	DeviceID   string
	Severity   string
	TempC      float64
	OccurredAt string
	AlertID    string
	Tags       []string
}

const chunkSize = 500

func buildThermalMessage(token string,  p ThermalAlertPayload)*messaging.Message{
	data := map[string]string{
		"deviceId": p.DeviceID,
		"eventId": p.AlertID,
		"occuredAt": p.OccurredAt,
	}

	if p.Severity != "" {
		data["severity"] = p.Severity
	}

	if p.TempC != 0 {
		data["tempC"] = fmt.Sprintf("%.1f", p.TempC)
	}

	var formattedTime string

	t, err := time.Parse(time.RFC3339, p.OccurredAt)
	if err == nil {
		formattedTime = t.Format("2006/01/02 15:04")
	} else {
		formattedTime = "取得に失敗しました。"
	}

	return &messaging.Message{
		Token: token,
		Data: data,
		Notification: &messaging.Notification{
			Title:   "異常な発熱を検知しました",
			Body:    fmt.Sprintf("デバイス：%s\n発生日時：%s", p.DeviceID, formattedTime),
		},
		//TODO colapseKeyとTagの設定の調整
		Android: &messaging.AndroidConfig{
			Priority: "high",
			CollapseKey: p.DeviceID,
			Notification: &messaging.AndroidNotification{
				ChannelID:     "alerts",
				ClickAction:   "OPEN_ALERT",
				Tag:           p.DeviceID,
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{"apns-priority": "10"},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Category: "ALERT",
					ThreadID: p.DeviceID,
					Sound:    "default",
				},
			},
		},
	}
}

func (f *FCM) SendThermalAlert(ctx context.Context, db *gorm.DB, p ThermalAlertPayload, tokens []string )error{
	for i := 0; i < len(tokens); i += chunkSize{
		end := i + chunkSize
		if end > len(tokens){
			end = len(tokens)
		}
		var msgs []*messaging.Message
		for _, t := range tokens[1:end]{
			msgs = append(msgs, buildThermalMessage(t, p))
		}

		resps, err := f.Client.SendEach(ctx, msgs)
		if err != nil {
			//TODO 再送処理
			return  errors.New("メッセージの送信に失敗しました：" + err.Error())
		}

		for idx, res := range resps.Responses{
			if !res.Success && res.Error != nil {
				if isInvalidToken(res.Error.Error()){
					_ = db.WithContext(ctx).Table("push_tokens").
					Where("token = ?", tokens[i+idx]).
					Updates(map[string]interface{}{
						"revoked": true,
						"revoked_reason": "auto："+ res.Error.Error(),
						"updated_at": time.Now(),
					}).Error
				}
			}
		}
	}
	return nil
}