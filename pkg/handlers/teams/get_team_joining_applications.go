package teams

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamJoiningApplications(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userTeamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	var applications []struct {
		ID                  uint
		CreatedAt           time.Time
		IntroductoryMessage string
		Candidate           string
	}

	h.DB.Raw(
		`SELECT tja.id, tja.created_at, tja.introductory_message, u.user_name AS candidate
		FROM team_joining_applications AS tja
		INNER JOIN users AS u ON u.id = tja.candidate_id
		WHERE tja.team_id = ?`, userTeamID,
	).Scan(&applications)

	if len(applications) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено заявок на вступление в вашу команду"})
		return
	}

	c.JSON(200, &applications)
}
