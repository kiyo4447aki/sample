package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "テスト用DBの作成に失敗しました")
	
	err = db.AutoMigrate(&User{}, &Device{}, &UserDevice{})
	require.NoError(t, err, "テスト用DBのマイグレーションに失敗しました")
	return db
}