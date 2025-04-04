package volumes

import (
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		INNER JOIN users ON user_roles.user_id = users.id
		WHERE users.id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь главой команды"})
		return
	}

	title := c.Param("title")
	volume := c.Param("volume")

	var volumeID uint
	h.DB.Raw(
		`SELECT volumes.id FROM volumes
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE titles.name = ? AND volumes.name = ?`,
		title, volume,
	).Scan(&volumeID)

	if volumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	h.DB.Raw(
		`SELECT
		(SELECT titles.team_id FROM titles
		INNER JOIN volumes ON titles.id = volumes.title_id
		WHERE volumes.id = ?)
		= 
		(SELECT team_id FROM users WHERE id = ?)`,
		volumeID, claims.ID,
	).Scan(&doesUserTeamTranslatesDesiredTitle)

	if !doesUserTeamTranslatesDesiredTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит тайтл, которому принадлежит этот том"})
		return
	}

	var volumeChapterID uint
	h.DB.Raw(
		`SELECT chapters.id FROM chapters AS c
		INNER JOIN volumes AS v ON v.id = c.volume_id 
		WHERE v.id = ?
		LIMIT 1`, volumeID,
	).Scan(&volumeChapterID)

	if volumeChapterID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "удалить можно только том без глав"}) // Тут можно было бы сделать каскадное удаление, но проблема в страницах глав, хранящихся в mongooDB. Они бы никуда не делись. Да и давать пользователю возможность одной кнопкой удалить том и все его главы идея тоже не лучшая.
		return
	}

	if result := h.DB.Exec("DELETE FROM volumes WHERE id = ?", volumeID); result.Error != nil {
		log.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error})
		return
	}

	c.JSON(200, gin.H{"success": "том успешно удалён"})
}
