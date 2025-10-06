package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//CreateDevice関数
//正常系
func TestCreateDevice_Success(t *testing.T){
	db := setupTestDB(t)

	device := &Device{
		DeviceID: "C12345",
		JanusPassword: "janus_pass",
		DeviceName: "testdevice",
		Location: "entrance",
	}

	err := CreateDevice(db, device)
	require.NoError(t, err, "CreateDeviceがエラーを返しました")

	//DBへ保存されているか確認
	var fetched Device
	err = db.First(&fetched, "device_id = ?", device.DeviceID).Error
	require.NoError(t, err, "DBからデバイスを取得できませんでした")
	assert.Equal(t, device.DeviceName, fetched.DeviceName, "Devicenameが期待値と異なります" )
	assert.Equal(t, device.Location, fetched.Location, "Devicenameが期待値と異なります" )
	assert.Equal(t, device.JanusPassword, fetched.JanusPassword, "Devicenameが期待値と異なります" )
	assert.NotZero(t,fetched.CreatedAt, "CreatedAtにnilが返されました")
	assert.NotZero(t, fetched.UpdatedAt, "UpdatedAtにnilが返されました")

}

//異常系：デバイスIDの重複
func TestCreateDevice_Duplicate(t *testing.T){
	db := setupTestDB(t)
	device := &Device{DeviceID: "dupID", JanusPassword: "pw"}
	err := CreateDevice(db, device)
	require.NoError(t, err, "初回 CreateDevice でエラーが発生しました")

	err = CreateDevice(db, device)
	require.Error(t, err, "重複した DeviceID でエラーが発生するはずです")
	assert.True(t,
		strings.Contains(err.Error(), "デバイスの登録に失敗しました"),
		"エラーメッセージが期待値と異なります: %s", err.Error(),
	)
}

//GetDeviceByID
//正常系
func TestGetDeviceByID_Success(t *testing.T) {
	db := setupTestDB(t)

	device := &Device{
		DeviceID: "C12345",
		JanusPassword: "janus_pass",
		DeviceName: "testdevice",
		Location: "entrance",
	}

	require.NoError(t, db.Create(device).Error, "テスト用デバイスの作成に失敗しました")

	fetched, err := GetDeviceByID(db, "C12345")
	require.NoError(t, err, "GetDeviceByID がエラーを返しました")
	assert.Equal(t, device.DeviceName, fetched.DeviceName, "DeviceName が一致しません")
	assert.Equal(t, device.Location, fetched.Location, "Location が一致しません")	
	assert.Equal(t, device.JanusPassword, fetched.JanusPassword, "Devicenameが一致しません" )
	assert.NotZero(t,fetched.CreatedAt, "CreatedAtにnilが返されました")
	assert.NotZero(t, fetched.UpdatedAt, "UpdatedAtにnilが返されました")
}

//異常系：デバイスが存在しない
func TestGetDeviceByID_NotFound(t *testing.T) {
	db := setupTestDB(t)

	_, err := GetDeviceByID(db, "noSuchID")
	require.Error(t, err, "存在しない DeviceID でエラーが発生するはずです")
	assert.Equal(t, "指定されたデバイスが見つかりません", err.Error(), "エラーメッセージが期待値と異なります")
}
