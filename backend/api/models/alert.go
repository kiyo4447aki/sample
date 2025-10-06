package models

import "time"

type Alert struct {
	ID         uint64    `gorm:"primaryKey"`
	DeviceID   string    `gorm:"type:varchar(64);not null;index:idx_alerts_device_occurred"`
	EventID    string    `gorm:"type:varchar(128);not null;uniqueIndex:ux_alerts_device_event"`
	Severity   string    `gorm:"type:varchar(16);not null;default:info"` 
	TempC      *float64  
	OccurredAt time.Time `gorm:"type:timestamptz;not null;index:idx_alerts_device_occurred,sort:desc"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Media AlertMedia `gorm:"foreignKey:AlertID"`
}
