package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T){
	secret := "testsecret"
	require.NoError(t, os.Setenv("JWT_SECRET", secret))
	require.NoError(t, os.Setenv("POSTGRES_DSN", "dummy"))
	require.NoError(t, os.Setenv("JANUS_WS", "dummy"))
}

func unsetEnv(){
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("POSTGRES_DSN")
	os.Unsetenv("JANUS_WS")
}

//GenerateJWTとValidateJWT関数
//正常系：トークンを生成し、検証できる
func TestGenerateAndValidateToken_Success(t *testing.T) {
	setEnv(t)
	defer unsetEnv()

	userID := 42
	role := "user"
	now := time.Now()

	tokenString, err := GenerateJWT(userID, role)
	require.NoError(t, err, "generateTokenがエラーを返しました")
	require.NotEmpty(t, tokenString, "トークンに空文字列が返されました")

	claims, err := ValidateJWT(tokenString)
	require.NoError(t, err, "ValidateJWTがエラーを返しました")

	assert.Equal(t, userID, claims.UserID, "UserID が期待値と異なります")
	assert.Equal(t, role, claims.Role, "Role が期待値と異なります")

	assert.WithinDuration(t, now, claims.IssuedAt.Time, time.Second, "IssuedAt が現在時刻と大きくずれています")

	expectedExpiry := now.Add(30 * 24 * time.Hour)
	assert.WithinDuration(t, expectedExpiry, claims.ExpiresAt.Time, time.Minute, "ExpiresAt が30日後と大きくずれています")
}

//異常系：秘密鍵未設定のとき、エラー発生
func TestGenerateToken_NoSecret(t *testing.T) {
	setEnv(t)
	defer unsetEnv()

	os.Unsetenv("JWT_SECRET")

	_, err := GenerateJWT(1, "user")
	require.Error(t, err, "秘密鍵未設定時にエラーが返されませんでした")
	assert.Equal(t, err.Error(), "環境変数 JWT_SECRET がセットされていません", "エラーメッセージが期待値と異なります")
}

// 異常系：不正な署名で生成したトークンは検証に失敗
func TestValidateToken_InvalidSignature(t *testing.T) {
	setEnv(t)
	defer unsetEnv()

	exp := time.Now().Add(10 * time.Minute)
	claims := &JwtClaims{
		UserID: 1,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "1",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	badToken, err := token.SignedString([]byte("badsecret"))
	require.NoError(t, err, "モックトークン(不正な秘密鍵使用)の生成に失敗しました")

	_, err = ValidateJWT(badToken)
	require.Error(t, err, "不正な署名のトークンでエラーが発生するはずです")
	assert.Contains(t, err.Error(), "JWTの解析に失敗しました", "エラーメッセージが期待値と異なります")
}
