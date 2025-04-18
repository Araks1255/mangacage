package joinrequests

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTeamJoinRequests(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var requests []struct {
		ID                  uint
		CreatedAt           time.Time
		IntroductoryMessage string
		Role                string
		Team                string
	}

	h.DB.Raw(
		`SELECT tjr.id, tjr.created_at, tjr.introductory_message, tjr.role, t.name AS team
		FROM team_join_requests AS tjr
		INNER JOIN teams AS t ON t.id = tjr.team_id
		WHERE tjr.candidate_id = ?`, claims.ID,
	).Scan(&requests)

	if len(requests) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет заявок на вступление в команду перевода"})
		return
	}

	c.JSON(200, &requests)
}
