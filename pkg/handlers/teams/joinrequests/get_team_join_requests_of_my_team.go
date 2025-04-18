package joinrequests

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamJoinRequestsOfMyTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var teamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	var requests []struct {
		ID                  uint
		CreatedAt           time.Time
		IntroductoryMessage string
		Role                string
		Candidate           string
	}

	h.DB.Raw(
		`SELECT tjr.id, tjr.created_at, tjr.introductory_message, tjr.role, u.user_name AS candidate
		FROM team_join_requests AS tjr
		INNER JOIN users AS u ON u.id = tjr.candidate_id
		WHERE tjr.team_id = ?`, teamID,
	).Scan(&requests)

	if len(requests) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено заявок на вступление в вашу команду"})
		return
	}

	c.JSON(200, &requests)
}
