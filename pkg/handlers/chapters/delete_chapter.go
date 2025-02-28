package chapters

import (
	"log"
	"slices"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles ON roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if isUserTeamLeader := slices.Contains(userRoles, "team_leader"); !isUserTeamLeader {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы не являетесь лидером команды перевода"})
		return
	}

	desiredChapter := strings.ToLower(c.Param("chapter"))

	var WasDesiredChapterCreatedByUserTeam bool
	h.DB.Raw("SELECT CAST(CASE WHEN (SELECT team_id FROM users WHERE id = ?) = (SELECT titles.team_id FROM titles INNER JOIN chapters ON titles.id = chapters.title_id WHERE chapters.name = ?) THEN TRUE ELSE FALSE END AS BOOLEAN)", claims.ID, desiredChapter).Scan(&WasDesiredChapterCreatedByUserTeam)

	if !WasDesiredChapterCreatedByUserTeam {
		c.AbortWithStatusJSON(403, gin.H{"error": "Удалить главу может только команда, выложившая её"})
		return
	}

	if result := h.DB.Exec("DELETE FROM chapters WHERE name = ?", desiredChapter); result.RowsAffected == 0 {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось удалить главу. Возможно, была допущена опечатка"})
		return
	}

	c.JSON(200, gin.H{"success": "Глава успешно удалена"})
}
