package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//AddUserDevice関数
//正常系
func TestAddUserDevice_Success(t *testing.T) {
	db := setupTestDB(t)

	user, err := CreateUser(db, "C123456", "passhash", "user")
	require.NoError(t, err, "ユーザー作成に失敗しました")
	device := &Device{DeviceID: "dev1"}
	require.NoError(t, err, "デバイス作成に失敗しました")
	err = AddUserDevice(db, user.ID, device.DeviceID)
	require.NoError(t, err, "AddUserDeviceがエラーを返しました")

	var uds []UserDevice
	require.NoError(t, db.Where("user_id = ?", user.ID).Find(&uds).Error, "DBからレコードを取得できませんでした")
	assert.Len(t, uds, 1, "レコード数が期待された値と異なります")
	assert.Equal(t, user.ID, uds[0].UserID, "UserID が一致しません")
	assert.Equal(t, device.DeviceID, uds[0].DeviceID, "DeviceID が一致しません")
}

//異常系：同じユーザー・デバイスの組み合わせを2回登録すると、重複エラー発生
func TestAddUserDevice_Duplicate(t *testing.T) {
	db := setupTestDB(t)

	user, err := CreateUser(db, "C123456", "pwhash", "user")
	require.NoError(t, err, "ユーザー作成に失敗しました")
	device := &Device{DeviceID: "dev2"}
	require.NoError(t, CreateDevice(db, device), "デバイス作成に失敗しました")

	require.NoError(t, AddUserDevice(db, user.ID, device.DeviceID))

	err = AddUserDevice(db, user.ID, device.DeviceID)
	require.Error(t, err, "重複した関連付けでエラーが発生しませんでした")
	assert.Equal(t, "既にこのユーザーとデバイスの関連付けは存在しています", err.Error(), "エラーメッセージが期待値と異なります")
}

//GetUserDevices関数
//正常系
func TestGetUserDevices_Success(t *testing.T) {
	db := setupTestDB(t)

	user, err := CreateUser(db, "C123456", "pwhash", "user")
	require.NoError(t, err)

	devA := &Device{DeviceID: "devA"}
	devB := &Device{DeviceID: "devB"}
	require.NoError(t, CreateDevice(db, devA))
	require.NoError(t, CreateDevice(db, devB))

	require.NoError(t, AddUserDevice(db, user.ID, devA.DeviceID))
	require.NoError(t, AddUserDevice(db, user.ID, devB.DeviceID))

	uds, err := GetUserDevices(db, user.ID)
	require.NoError(t, err, "GetUserDevices がエラーを返しました")
	assert.Len(t, uds, 2, "関連付けの件数が期待と異なります")

	ids := []string{uds[0].DeviceID, uds[1].DeviceID}
	assert.Contains(t, ids, devA.DeviceID, "devA が含まれていません")
	assert.Contains(t, ids, devB.DeviceID, "devB が含まれていません")
}

// 正常系：関連付けが一切ない場合、空スライスが返る
func TestGetUserDevices_Empty(t *testing.T) {
	db := setupTestDB(t)

	user, err := CreateUser(db, "C123456", "pwhash", "user")
	require.NoError(t, err)

	uds, err := GetUserDevices(db, user.ID)
	require.NoError(t, err, "GetUserDevices がエラーを返しました")
	assert.Empty(t, uds, "関連付けが存在しないはずですが、結果が空ではありません")
}

