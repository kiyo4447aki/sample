package middleware

import (
	"net/http"
	"strings"

	"backend-proto/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//JWTを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, 
				gin.H{
					"status": "failed",
					"error": "認証ヘッダーが存在しません"},
			)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "failed",
				"error":  "認証ヘッダーが不正です",
			})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))


		//JWTからclaimsを取得
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"status": "failed",
					"error": err.Error()},
			)
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}