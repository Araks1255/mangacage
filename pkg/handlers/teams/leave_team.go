package teams

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) LeaveTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var user models.User
	h.DB.Raw("SELECT * FROM users WHERE id = ?", claims.ID).Scan(&user)

	if !user.TeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы итак не состоите в команде перевода"})
		return
	}

	user.TeamID.Valid = false
	user.TeamID.Int32 = 0

	transaction := h.DB.Begin()

	if result := transaction.Save(&user); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось удалить вас из команды"})
		return
	}

	if result := transaction.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = (SELECT id FROM roles WHERE name = 'translater')", claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось снять вас с роли переводчика"})
		return
	}

	transaction.Commit()

	c.JSON(200, gin.H{"success": "Вы успешно покинули команду перевода"})
}
