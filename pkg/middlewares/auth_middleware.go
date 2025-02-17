package middlewares

import (
	"log"

	"github.com/Araks1255/mangabrad/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("mangabrad_token")
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Вы не авторизованы"})
			return
		}

		claims, err := utils.ParseToken(cookie, secretKey)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Вы не авторизованы"})
			return
		}

		c.Set("claims", claims)
	}
}
