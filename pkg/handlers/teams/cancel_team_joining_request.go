package teams

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelTeamJoiningRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	team := c.Param("team")

	var teamID uint
	h.DB.Raw("SELECT id FROM teams WHERE name = ?", team).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	var userTeamJoiningRequestID uint
	h.DB.Raw("SELECT id FROM team_joining_applications WHERE candidate_id = ? AND team_id = ?", claims.ID, teamID).Scan(&userTeamJoiningRequestID)
	if userTeamJoiningRequestID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет заявки на вступление в эту команду"})
		return
	}

	if result := h.DB.Exec("DELETE FROM team_joining_applications WHERE id = ?", userTeamJoiningRequestID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "ваша заявка на вступление в команду успешно отменена"})
}
