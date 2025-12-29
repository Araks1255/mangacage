package joinrequests

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamJoinRequestsOfMyTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requests []dto.ResponseTeamJoinRequestDTO

	if err := h.DB.Raw(
		`SELECT
			tjr.id, tjr.candidate_id, u.user_name AS candidate
		FROM
			team_join_requests AS tjr
			INNER JOIN users AS u ON u.id = tjr.candidate_id
		WHERE
			tjr.team_id = (SELECT team_id FROM users WHERE id = ?)`,
		claims.ID,
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
