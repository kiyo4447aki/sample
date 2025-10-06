package models

import "time"

type AlertMedia struct {
	ID        uint64    `gorm:"primaryKey"`
	AlertID   *uint64   `gorm:"index"` 
	ObjectKey string    `gorm:"type:text;uniqueIndex;not null"`
	MimeType  string    `gorm:"type:text;not null"`
	SizeBytes *int64
	Width     *int
	Height    *int
	SHA256    *string   `gorm:"type:char(64)"`
	Committed bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
