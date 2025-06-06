package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTeamOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	moderationType := c.Query("type")

	var (
		team models.TeamOnModerationDTO
		err  error
	)

	switch moderationType {
	case "new":
		err = h.DB.Raw(
			"SELECT id, created_at, name, description FROM teams_on_moderation WHERE creator_id = ?", claims.ID,
		).Scan(&team).Error

	case "edited":
		err = h.DB.Raw(
			`SELECT
				tom.id, tom.created_at, tom.name, tom.description,
				t.name AS existing, t.id AS existing_id,
				u.user_name AS leader, u.id AS leader_id
			FROM
				teams_on_moderation AS tom
				INNER JOIN teams AS t ON t.id = tom.existing_id
				INNER JOIN users AS u ON u.team_id = t.id
				INNER JOIN user_roles AS ur ON u.id = ur.user_id
				INNER JOIN roles AS r ON r.id = ur.role_id
			WHERE
				tom.creator_id = ?
			AND
				r.name = 'team_leader'
			LIMIT 1`,
			claims.ID,
		).Scan(&team).Error

	case "":
		err = h.DB.Raw(
			`SELECT
				tom.id, tom.created_at, tom.name, tom.description,
				t.name AS existing, t.id AS existing_id,
				u.user_name AS leader, u.id AS leader_id
			FROM
				teams_on_moderation AS tom
				LEFT JOIN teams AS t ON tom.existing_id = t.id
				LEFT JOIN users AS u ON t.id = u.team_id
				LEFT JOIN user_roles AS ur ON u.id = ur.user_id AND ur.role_id = (SELECT id FROM roles WHERE name = 'team_leader')
			WHERE
				tom.creator_id = ?
			LIMIT 1`,
			claims.ID,
		).Scan(&team).Error

	default:
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный тип модерации"})
		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if team.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено вашей заявки на модерацию команды"})
		return
	}

	c.JSON(200, &team)
}
