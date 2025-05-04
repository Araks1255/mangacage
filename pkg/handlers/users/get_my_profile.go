package users

import (
	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var user models.UserDTO

	h.DB.Raw(
		`SELECT u.id, u.created_at, u.user_name, u.about_yourself,
		t.name AS team, t.id AS team_id,
		(
			SELECT ARRAY(
				SELECT r.name FROM roles AS r
				INNER JOIN user_roles AS ur ON ur.role_id = r.id
				WHERE ur.user_id = u.id
			) AS roles
		) FROM users AS u
		LEFT JOIN teams AS t ON u.id = t.id
		WHERE u.id = ?`, claims.ID,
	).Scan(&user)

	if user.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "произошла ошибка при получении профиля"})
		return
	}

	c.JSON(200, &user)
}
