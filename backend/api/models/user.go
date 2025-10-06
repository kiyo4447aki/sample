package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	HashedPassword string    `gorm:"not null" json:"-"`
	Role	  string    `gorm:"not null;default:user" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateUser(db *gorm.DB, username, hashedPassword, role string)(*User, error){
	
	user := &User{
		Username:      username,
		HashedPassword: hashedPassword,
		Role:          role,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("ユーザーの作成に失敗しました: %w", err)
	}
	return user, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ユーザー %s は存在しません", username)
		}
		return nil, fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
	}
	return &user, nil
}