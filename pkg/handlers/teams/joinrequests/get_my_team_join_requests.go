package joinrequests

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTeamJoinRequests(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requests []dto.ResponseTeamJoinRequestDTO

	if err := h.DB.Raw(
		`SELECT
			tjr.id, tjr.team_id, t.name AS team
		FROM
			team_join_requests AS tjr
			INNER JOIN teams AS t ON t.id = tjr.team_id
		WHERE
			tjr.candidate_id = ?`,
		claims.ID,
	).Scan(&requests).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(requests) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет заявок на вступление в команду перевода"})
		return
	}

	c.JSON(200, &requests)
}
