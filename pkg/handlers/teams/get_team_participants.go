package teams

import (
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamParticipants(c *gin.Context) {
	team := c.Param("team")

	var teamID uint
	h.DB.Raw("SELECT id FROM teams WHERE name = ?", team).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
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
