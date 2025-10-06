package models

import (
	"errors"

	"gorm.io/gorm"
)

type UserDevice struct {
	ID       uint   `gorm:"primaryKey" json:"id"`                         
	UserID   uint   `gorm:"not null" json:"user_id"`                       
	DeviceID string `gorm:"not null;type:varchar(100)" json:"device_id"`   
}

func AddUserDevice(db *gorm.DB, userID uint, deviceID string) error {
	var count int64
	if err := db.Model(&UserDevice{}).Where("user_id = ? AND device_id = ?", userID, deviceID).Count(&count).Error; err != nil {
		return errors.New("ユーザーとデバイスの関連付けの状態確認に失敗しました: " + err.Error())
	}
	if count > 0 {
		return errors.New("既にこのユーザーとデバイスの関連付けは存在しています")
	}

	userDevice := &UserDevice{
		UserID:   userID,
		DeviceID: deviceID,
	}
	if err := db.Create(userDevice).Error; err != nil {
		return errors.New("ユーザーとデバイスの関連付けの登録に失敗しました: " + err.Error())
	}
	return nil
}

func GetUserDevices(db *gorm.DB, userID uint) ([]UserDevice, error) {
	var userDevices []UserDevice
	if err := db.Where("user_id = ?", userID).Find(&userDevices).Error; err != nil {
		return nil, errors.New("ユーザーに紐付いたデバイス情報の取得に失敗しました: " + err.Error())
	}
	return userDevices, nil
}
