package teams

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id команды должен быть числом"})
		return
	}

	var team models.TeamDTO

	if err = h.DB.Raw(
		`SELECT
			t.id, t.created_at, t.name, t.description,
			u.user_name AS leader, u.id AS leader_id
		FROM
			teams AS t
			INNER JOIN users AS u ON u.team_id = t.id
			INNER JOIN user_roles AS ur ON ur.user_id = u.id
			INNER JOIN roles AS r ON r.id = ur.role_id
		WHERE
			r.name = 'team_leader'
		AND
			t.id = ?`,
		teamID,
	).Scan(&team).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if team.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	c.JSON(200, &team)
}
