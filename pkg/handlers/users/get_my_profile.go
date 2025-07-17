package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var user dto.ResponseUserDTO

	err := h.DB.Table("users AS u").
		Select(
			`u.*, t.name AS team, ARRAY(
				SELECT r.name FROM roles AS r
				INNER JOIN user_roles AS ur ON ur.role_id = r.id
				WHERE ur.user_id = u.id
			) AS roles`,
		).
		Joins("LEFT JOIN teams AS t ON u.team_id = t.id").
		Where("u.id = ?", claims.ID).
		Scan(&user).Error

	if err != nil {
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
