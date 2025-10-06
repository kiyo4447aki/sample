package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//TestCreateUser関数
//正常系
func TestCreateUser_Success(t *testing.T){
	db := setupTestDB(t)
	user, err := CreateUser(db, "username", "hashedpw", "user")
	require.NoError(t, err, "CreateUserがエラーを返しました")
	require.NotZero(t, user.ID, "ユーザーIDにゼロ値が返されました")
	assert.Equal(t, "username", user.Username, "Usernameが期待される値と異なります")
	assert.Equal(t, "hashedpw", user.HashedPassword, "HashedPasswordが期待される値と異なります")
	assert.Equal(t, "user", user.Role, "が期待される値と異なります")
	assert.NotZero(t,user.CreatedAt, "CreatedAtにnilが返されました")
	assert.NotZero(t,user.UpdatedAt, "UpdatedAtにnilが返されました")
	
}

//異常系：ユーザーIDの重複
func TestCreateUser_Duplicate(t *testing.T){
	db := setupTestDB(t)
	_, err := CreateUser(db, "username", "hashedpw", "user")
	require.NoError(t, err, "初回のCreateUserでエラーが発生しました")

	_, err = CreateUser(db, "username", "hashedpw", "user")
	require.Error(t, err, "usernameの重複時、エラーが発生するはずですが、エラーが返されませんでした")
}

//TestGetUserByUsername関数
//正常系
func TestGetUserByUsername_success(t *testing.T){
	db := setupTestDB(t)
	_, err := CreateUser(db, "username", "hashedpw", "user")
	require.NoError(t, err, "テスト用ユーザーの作成でエラーが発生しました")

	user, err := GetUserByUsername(db, "username")
	require.NoError(t, err, "ユーザーの取得時にエラーが発生しました")
	require.NotZero(t, user.ID, "ユーザーIDにゼロ値が返されました")
	assert.Equal(t, "username", user.Username, "Usernameが期待される値と異なります")
	assert.Equal(t, "hashedpw", user.HashedPassword, "HashedPasswordが期待される値と異なります")
	assert.Equal(t, "user", user.Role, "が期待される値と異なります")
	assert.NotZero(t,user.CreatedAt, "CreatedAtにnilが返されました")
	assert.NotZero(t,user.UpdatedAt, "UpdatedAtにnilが返されました")
}

//異常系：ユーザーが存在しない
func TestGetUserByUsername_NotFound(t *testing.T){
	db := setupTestDB(t)
	_, err := GetUserByUsername(db, "username")
	require.Error(t, err, "存在しないユーザーを検索した際、エラーが返されるはずですが、エラーが発生しませんでした")
	assert.Equal(t, "ユーザー username は存在しません", err.Error(), "期待されるエラーと異なるエラーが返されました")
}
