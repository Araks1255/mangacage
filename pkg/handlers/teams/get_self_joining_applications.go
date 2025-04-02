package teams

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSelfJoiningApplications(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var applications []struct {
		ID                  uint
		CreatedAt           time.Time
		IntroductoryMessage string
		Team                string
	}

	h.DB.Raw(
		`SELECT tja.id, tja.created_at, tja.introductory_message, t.name AS team
		FROM team_joining_applications AS tja
		INNER JOIN teams AS t ON t.id = tja.team_id
		WHERE tja.candidate_id = ?`, claims.ID,
	).Scan(&applications)

	if len(applications) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет заявок на вступление в команду перевода"})
		return
	}

	c.JSON(200, &applications)
}
