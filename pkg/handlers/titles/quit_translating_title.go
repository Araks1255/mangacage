package titles

import (
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"github.com/gin-gonic/gin"
)

func (h handler) QuitTranslatingTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON user_roles.role_id = roles.id
		WHERE user_roles.user_id = ?`, claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	title := c.Param("title")

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE lower(name) = lower(?)", title).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var IsUserTeamTranslatesThisTitle bool
	h.DB.Raw(
		`SELECT (select team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)`,
		titleID, claims.ID,
	).Scan(&IsUserTeamTranslatesThisTitle)

	if !IsUserTeamTranslatesThisTitle && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	if result := h.DB.Exec("UPDATE titles SET team_id = null WHERE id = ?", titleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "ваша команда больше не переводит данный тайтл"})
}
