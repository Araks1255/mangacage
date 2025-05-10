package middlewares

import (
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequireRoles(db *gorm.DB, roles []string) gin.HandlerFunc { // Я вообще хотел вынести роли в jwt токен, но так подумал, что не хочу, чтобы пользователь, которого сняли с роли заместителя лидера команды мог продолжать пользоваться своими привелегиями до рефреша токена. С бд всё-таки гораздо безопаснее.
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*auth.Claims)

		var userRoles []string

		if err := db.Raw(
			`SELECT
				r.name
			FROM
				roles AS r
				INNER JOIN user_roles AS ur ON ur.role_id = r.id
			WHERE
				ur.user_id = ?`,
			claims.ID,
		).Scan(&userRoles).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		var accessAllowed bool

		for i := 0; i < len(roles); i++ {
			if slices.Contains(userRoles, roles[i]) {
				accessAllowed = true
				break
			}
		}

		if !accessAllowed {
			c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для совершения этого действия"})
			return
		}

		c.Next()
	}
}
