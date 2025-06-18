package middlewares

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func RequireRoles(db *gorm.DB, roles []string) gin.HandlerFunc { // Я вообще хотел вынести роли в jwt токен, но так подумал, что не хочу, чтобы пользователь, которого сняли с роли заместителя лидера команды мог продолжать пользоваться своими привелегиями до рефреша токена. С бд всё-таки гораздо безопаснее.
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*auth.Claims)

		var allowed bool

		err := db.Raw(
			`SELECT EXISTS(
				SELECT 1 FROM user_roles AS ur
				INNER JOIN roles AS r ON r.id = ur.role_id
				WHERE ur.user_id = ?
				AND r.name = ANY(?::TEXT[])
			)`, claims.ID, pq.Array(roles),
		).Scan(&allowed).Error

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для совершения этого действия"})
			return
		}

		c.Next()
	}
}
