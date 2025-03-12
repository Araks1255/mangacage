package titles

import (
	"database/sql"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"github.com/gin-gonic/gin"
)

func (h handler) TranslateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles on roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if IsUserTeamOwner := slices.Contains(userRoles, "team_leader"); !IsUserTeamOwner {
		c.AbortWithStatusJSON(403, gin.H{"error": "Взять тайтл на перевод может только владелец команды перевода"})
		return
	}

	title := c.Param("title")

	var desiredTitle models.Title
	h.DB.Raw("SELECT * FROM titles WHERE lower(name) = lower(?)", title).Scan(&desiredTitle)
	if desiredTitle.TeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "Тайтл уже переводит другая команда"})
		return
	}

	if desiredTitle.OnModeration {
		c.AbortWithStatusJSON(403, gin.H{"error": "Этот тайтл находится на стадии модерации"})
		return
	}

	var userTeamID sql.NullInt64
	h.DB.Raw("SELECT teams.id FROM teams INNER JOIN users ON teams.id = users.team_id WHERE users.id = ?", claims.ID).Scan(&userTeamID)
	if !userTeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы не состоите в команде перевода"})
		return
	}

	desiredTitle.TeamID = userTeamID

	if result := h.DB.Save(&desiredTitle); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось взять тайтл на перевод"})
		return
	}

	c.JSON(200, gin.H{"success": "Теперь ваша команда переводит этот тайтл"})
}
