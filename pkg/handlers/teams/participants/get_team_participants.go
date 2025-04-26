package participants

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamParticipants(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id команды должен быть числом"})
		return
	}

	var participants []models.UserDTO

	h.DB.Raw(
		`SELECT u.id, u.user_name,
		ARRAY_AGG(r.name)::TEXT[] AS roles
		FROM users AS u
		LEFT JOIN user_roles AS ur ON u.id = ur.user_id
		LEFT JOIN roles AS r ON r.id = ur.role_id 
		WHERE u.team_id = ?
		GROUP BY u.id`,
		teamID,
	).Scan(&participants)

	if len(participants) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "в этой команде нет участников"})
		return
	}

	c.JSON(200, &participants)
}
