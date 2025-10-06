package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Device struct {
	DeviceID          string    `gorm:"primaryKey;type:varchar(100)" json:"device_id"` 
	JanusPassword     string    `gorm:"not null" json:"janus_password"`               
	DeviceName        string    `gorm:"type:varchar(255)" json:"device_name"`         
	Location          string    `gorm:"type:varchar(255)" json:"location"`            
	Status            string    `gorm:"type:varchar(50)" json:"status"`               
	LastCommunication time.Time `json:"last_communication"`                           
	CreatedAt         time.Time `json:"created_at"`                                   
	UpdatedAt         time.Time `json:"updated_at"`                                   
}


func CreateDevice(db *gorm.DB, device *Device) error {
	if err := db.Create(device).Error; err != nil {
		return errors.New("デバイスの登録に失敗しました: " + err.Error())
	}
	return nil
}

func GetDeviceByID(db *gorm.DB, deviceID string) (*Device, error) {
	var device Device
	if err := db.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("指定されたデバイスが見つかりません")
		}
		return nil, errors.New("デバイス情報の取得に失敗しました: " + err.Error())
	}
	return &device, nil
}

