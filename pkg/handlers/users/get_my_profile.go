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

	err := h.DB.Table("users AS u").Select("u.*, t.name AS team, ARRAY_AGG(r.name) AS roles").
		Joins("LEFT JOIN teams AS t ON u.team_id = t.id").
		Joins("LEFT JOIN user_roles AS ur ON u.id = ur.user_id").
		Joins("LEFT JOIN roles AS r ON ur.role_id = r.id").
		Where("u.id = ?", claims.ID).
		Group("u.id, t.id").
		Scan(&user).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	log.Printf("%+v", user)

	if user.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "произошла ошибка при получении профиля"})
		return
	}

	c.JSON(200, &user)
}
