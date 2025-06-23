package users

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id пользователя"})
		return
	}

	var user models.UserDTO

	err = h.DB.Table("users AS u").
		Select("u.*, ARRAY_AGG(DISTINCT r.name) AS roles, t.name AS team").
		Joins("INNER JOIN user_roles AS ur ON ur.user_id = u.id").
		Joins("INNER JOIN roles AS r ON r.id = ur.role_id").
		Joins("LEFT JOIN teams AS t ON u.team_id = t.id").
		Where("u.id = ?", userID).
		Where("u.visible").
		Group("u.id, t.id").
		Scan(&user).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if user.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "пользователь не найден"})
		return
	}

	c.JSON(200, &user)
}
