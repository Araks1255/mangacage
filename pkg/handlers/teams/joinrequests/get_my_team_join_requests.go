package joinrequests

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTeamJoinRequests(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requests []struct {
		ID                  uint      `json:"id"`
		CreatedAt           time.Time `json:"createdAt"`
		IntroductoryMessage string    `json:"introductoryMessage"`
		Role                string    `json:"role"`
		Team                string    `json:"team"`
		TeamID              uint      `json:"team_id"`
	}

	h.DB.Raw(
		`SELECT tjr.id, tjr.created_at, tjr.introductory_message, tjr.role,
		t.name AS team, t.id AS team_id
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
