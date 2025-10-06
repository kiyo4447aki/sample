package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"backend-proto/models"
	"backend-proto/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"` 
	Password string `json:"password" binding:"required"` 
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"` 
	Password string `json:"password" binding:"required"` 
}

const (
	headerAuthorization = "Authorization"
	headerDeviceID      = "X-Device-Id"
)

func RegisterUserHandler(c *gin.Context){
	//リクエストデータを取得
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, 
			gin.H{
				"status": "failed",
				"error":"リクエストデータが不正です"+err.Error(),
			},
		)
		return
	}

	//パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error": "ユーザー登録に失敗しました"+err.Error(),
		})
		return
	}

	//ユーザーの作成
	//roleは"user"をデフォルトとする
	defaultRole := "user"
	user, err := models.CreateUser(db, req.Username, string(hashedPassword), defaultRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error": "ユーザー登録に失敗しました"+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user_id": user.ID,
	})

}

func LoginHandler(c *gin.Context){
	//リクエストデータを取得
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, 
			gin.H{
				"status": "failed",
				"error":"リクエストデータが不正です"+err.Error(),
			},
		)
		return
	}

	//ユーザー情報の取得
	user, err := models.GetUserByUsername(db, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, 
			gin.H{
				"status": "failed",
				"error":"ユーザー名またはパスワードが正しくありません",
			},
		)
		return
	}

	//パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, 
			gin.H{
				"status": "failed",
				"error":"ユーザー名またはパスワードが正しくありません",
			},
		)
		return
	}

	//JWTトークンの生成
	token, err := utils.GenerateJWT(int(user.ID), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error": "トークン生成に失敗しました"+err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token": token,
	})
}

func TokenRefreshHandler(c *gin.Context){
	//claimsを取得
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error": "認証情報が見つかりません",
		})
		return
	}

	claims, ok := claimsInterface.(*utils.JwtClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "failed",
			"error": "認証情報の形式が正しくありません",
		})
		return
	}

	//新しいトークンを生成
	newToken, err := utils.GenerateJWT(claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error": "トークンの更新に失敗しました" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

func AuthRequestHandler(c *gin.Context) {
	// Authorization ヘッダから Bearer トークンを取り出す
	authz := c.GetHeader(headerAuthorization)
	if authz == "" {
		c.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}
	tokenStr := strings.TrimSpace(authz)
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	}
	if tokenStr == "" {
		c.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	// JWT の検証
	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil || claims == nil {
		// 署名不正・期限切れなど
		c.AbortWithStatus(http.StatusUnauthorized) // 401
		return
	}

	// デバイス認可チェック
	deviceID := c.GetHeader(headerDeviceID)
	if deviceID != "" {
		ok := hasDeviceAccess(claims.UserID, deviceID)
		if !ok {
			c.AbortWithStatus(http.StatusForbidden) // 403
			return
		}
	}

	c.Header("X-User-Id",  intToString(claims.UserID))
	c.Header("X-User-Role", claims.Role)
	c.Status(http.StatusNoContent) // 204
}

func intToString(v int) string {
	return fmt.Sprintf("%d", v)
}
