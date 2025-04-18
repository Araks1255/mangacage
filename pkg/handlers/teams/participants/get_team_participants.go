package participants

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamParticipants(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error":"id команды должен быть числом"})
		return
	}

	var participants []struct {
		ID       uint
		UserName string
		Role     string
	}

	h.DB.Raw(
		`SELECT u.id, u.user_name, r.name AS role FROM users AS u
		INNER JOIN user_roles AS ur ON u.id = ur.user_id
		INNER JOIN roles AS r ON r.id = ur.role_id
		INNER JOIN teams AS t ON t.id = u.team_id
		WHERE t.id = ? AND r.type = 'team'`,
		teamID,
	).Scan(&participants)

	if len(participants) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "в этой команде нет участников"})
		return
	}

	c.JSON(200, &participants)
}
