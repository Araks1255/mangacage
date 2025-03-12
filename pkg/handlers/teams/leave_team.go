package teams

import (
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) LeaveTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userTeamID sql.NullInt64
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)

	if !userTeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы итак не состоите в команде перевода"})
		return
	}

	tx := h.DB.Begin()

	if result := tx.Exec("UPDATE users SET team_id = null WHERE id = ?", claims.ID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result := tx.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = (SELECT id FROM roles WHERE name = 'translater')", claims.ID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "Вы успешно покинули команду перевода"})
}
