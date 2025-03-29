package titles

import (
	"database/sql"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")

	var (
		titleID     uint
		titleTeamID sql.NullInt64
	)

	row := h.DB.Raw("SELECT id, team_id FROM titles WHERE lower(name) = lower(?)", title).Row()
	row.Scan(&titleID, &titleTeamID)

	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": " тайтл не найден"})
		return
	}

	var titleVolumeID uint
	h.DB.Raw(`SELECT id FROM volumes WHERE title_id = ? LIMIT 1`, titleID).Scan(&titleVolumeID)
	if titleVolumeID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "удалить можно только тайтл без томов"})
		return
	}

	if !titleTeamID.Valid {
		var wasDesiredTitleCreatedByUser bool
		h.DB.Raw("SELECT (SELECT creator_id FROM titles WHERE id = ?) = ?", titleID, claims.ID).Scan(&wasDesiredTitleCreatedByUser)
		if !wasDesiredTitleCreatedByUser {
			c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь создателем этой главы"})
			return
		}

		if result := h.DB.Exec("DELETE FROM titles WHERE id = ? AND team_id IS NULL", titleID); result.Error != nil {
			log.Println(result.Error)
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
			return
		}
	}

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	h.DB.Raw("SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID).Scan(&doesUserTeamTranslatesDesiredTitle)
	if !doesUserTeamTranslatesDesiredTitle && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	// ПРОВЕРИТЬ ПОТОМ НИЧЕГО ЛИ НЕ ЗАБЫЛ

	if result := h.DB.Exec("DELETE FROM titles CASCADE WHERE id = ?", titleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "тайтл успешно удалён"})
}
