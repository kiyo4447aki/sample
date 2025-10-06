// utils/db.go
package utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	// GORM を用いてデータベースに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("GORM を用いたデータベース接続の初期化に失敗しました: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("内部データベースインスタンスの取得に失敗しました: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("データベースへの接続確認に失敗しました: %w", err)
	}

	return db, nil
}

