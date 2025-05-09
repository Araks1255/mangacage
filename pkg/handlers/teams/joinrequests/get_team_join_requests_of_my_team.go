package joinrequests

import (
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamJoinRequestsOfMyTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var teamID sql.NullInt64

	if err := h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&teamID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !teamID.Valid {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	var requests []models.TeamJoinRequestDTO

	if err := h.DB.Raw(
		`SELECT
			tjr.id, tjr.created_at, tjr.introductory_message,
			r.id AS role_id, r.name AS role,
			c.id AS candidate_id, c.user_name AS candidate
		FROM
			team_join_requests AS tjr
			LEFT JOIN roles AS r ON tjr.role_id = r.id
			INNER JOIN users AS c ON c.id = tjr.candidate_id
		WHERE
			tjr.team_id = ?`,
		teamID.Int64,
	).Scan(&requests).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(requests) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено заявок на вступление в вашу команду"})
		return
	}

	c.JSON(200, &requests)
}
