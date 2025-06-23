package middlewares

import (
	"github.com/Araks1255/mangacage/pkg/auth/utils"
	"github.com/gin-gonic/gin"
)

func Auth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.ParseClaims(c.Cookie, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "вы не авторизованы"})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
