package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var user models.UserDTO

	if err := h.DB.Raw(
		`SELECT
			u.id, u.created_at, u.user_name, u.about_yourself,
			t.name AS team, t.id AS team_id,
			ARRAY_AGG(r.name) AS roles
		FROM
			users AS u
			LEFT JOIN teams AS t ON u.team_id = t.id
			LEFT JOIN user_roles AS ur ON u.id = ur.user_id
			LEFT JOIN roles AS r ON ur.role_id = r.id
		WHERE
			u.id = ?
		GROUP BY
			u.id, t.id`,
		claims.ID,
	).Scan(&user).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if user.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "произошла ошибка при получении профиля"})
		return
	}

	c.JSON(200, &user)
}
