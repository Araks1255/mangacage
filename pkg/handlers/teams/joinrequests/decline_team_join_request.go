package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) DeclineTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки"})
		return
	}

	result := h.DB.Exec("DELETE FROM team_join_requests WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)", desiredRequestID, claims.ID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка на вступление в вашу команду не найдена"})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на вступление в вашу команду успешно отменена"})
	// Уведомление кандидату
}
