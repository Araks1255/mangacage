package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamJoinRequestOfMyTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	joinRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки на вступление в команду"})
		return
	}

	var result dto.ResponseTeamJoinRequestDTO

	err = h.DB.Raw(
		`SELECT
			tjr.*, u.user_name AS candidate, r.name AS role
		FROM
			team_join_requests AS tjr
			INNER JOIN users AS u ON u.id = tjr.candidate_id
			LEFT JOIN roles AS r ON tjr.role_id = r.id 
		WHERE
			tjr.id = ? AND tjr.team_id = (SELECT team_id FROM users WHERE id = ?)`,
		joinRequestID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка не найдена среди заявок на вступление в вашу команду"})
		return
	}

	c.JSON(200, &result)
}
