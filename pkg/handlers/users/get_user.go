package users

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id пользователя"})
		return
	}

	var user dto.ResponseUserDTO

	err = h.DB.Table("users AS u").
		Select(
			`u.*, t.name AS team, ARRAY(
				SELECT r.name FROM roles AS r
				INNER JOIN user_roles AS ur ON ur.role_id = r.id
				WHERE ur.user_id = u.id
			) AS roles`,
		).
		Joins("LEFT JOIN teams AS t ON u.team_id = t.id").
		Where("u.id = ?", userID).
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

	if !user.Visible {
		c.AbortWithStatusJSON(403, gin.H{"error": "этот пользователь имеет скрытый профиль"})
		return
	}

	c.JSON(200, &user)
}
