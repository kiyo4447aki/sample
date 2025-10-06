package models

import (
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PushToken struct {
	ID          uint64    `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null;index"`
	Token       string    `gorm:"type:varchar(2048);uniqueIndex;not null"`
	Platform    string    `gorm:"type:varchar(20);not null"` // "ios" | "android"
	LastSeenAt  time.Time `gorm:"autoUpdateTime"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Revoked     bool      `gorm:"default:false"`
	RevokedReason string  `gorm:"type:varchar(255)"`
}

func RegisterToken(db *gorm.DB ,tokenDB PushToken)error{
	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "token"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"user_id": tokenDB.UserID,
			"platform": tokenDB.Platform,
			"revoked": false,
			"revoked_reason": gorm.Expr("NULL"),
			"last_seen_at": time.Now(),
			"updated_at": time.Now(),
		}),
	}).Create(&tokenDB).Error; err!= nil {
		return errors.New("プッシュ通知トークンの登録に失敗しました：" + err.Error())
	}
	return nil
}

//ユーザーに紐づいたPushToken構造体を返す
func GetAvailableTokensByUserID(db *gorm.DB, userId string)([]PushToken, error){
	uid, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return nil, errors.New("ユーザーIDの形式が不正です：" + err.Error())
	}

	var tokens []PushToken
	if err := db.Where(
		"user_id = ? AND revoked = FALSE",
		uint(uid),
	).Order("updated_at DESC").Find(&tokens).Error; err != nil {
		return nil, errors.New("トークンを取得できませんでした" + err.Error())
	}
	return tokens, nil
}

//トークン文字列の返す
func GetAvailableTokensByDeviceID(db *gorm.DB, deviceId string)([]string, error){
	var tokens []string
	if err := db.Raw(`
		SELECT pt.token
		FROM user_devices ud
		JOIN push_tokens pt ON pt.user_id = ud.user_id
		WHERE ud.device_id = ? AND pt.revoked = FALSE
	`, deviceId).Scan(&tokens).Error; err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return []string{}, nil
	}
	return tokens, nil
}
