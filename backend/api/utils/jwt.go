package utils

import (
	"fmt"
	"time"

	"backend-proto/config"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserID int `json:"sub"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, role string)(string, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	jwtSecret := config.JWTSecret


	//有効期限を30日後に設定
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &JwtClaims{
		UserID: userID,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject: fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("JWTの生成に失敗しました: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*JwtClaims, error) {
	config, err := config.LoadConfig()
	jwtSecret := config.JWTSecret
	if err != nil {
		return nil, err
	}

	claims := &JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("JWTの解析に失敗しました: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("JWTが無効です")
	}
	return claims, nil
}
