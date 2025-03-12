package volumes

import (
	"fmt"
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
		WHERE users.id = ?`, claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь ни главой команды, ни модератором, ни администратором"})
		return
	}

	title := c.Param("title")
	volume := c.Param("volume")

	var volumeID uint
	h.DB.Raw(`SELECT volumes.id FROM volumes
	INNER JOIN titles ON volumes.title_id = titles.id
	WHERE lower(titles.name) = lower(?)
	AND lower(volumes.name) = lower(?)`,
		title,
		volume,
	).Scan(&volumeID)

	if volumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	h.DB.Raw(`SELECT CAST(
		CASE WHEN
		(SELECT titles.team_id FROM titles
		INNER JOIN volumes ON titles.id = volumes.title_id
		WHERE volumes.id = ?)
		= 
		(SELECT teams.id FROM teams
		INNER JOIN users ON teams.id = users.team_id
		WHERE users.id = ?)
		THEN TRUE ELSE FALSE END AS BOOLEAN)`,
		volumeID, claims.ID).
		Scan(&doesUserTeamTranslatesDesiredTitle)

	if !doesUserTeamTranslatesDesiredTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит тайтл, которому принадлежит этот том"})
		return
	}

	var volumeChapters []string
	h.DB.Raw("SELECT name FROM chapters WHERE volume_id = ?", volumeID).Scan(&volumeChapters)

	if len(volumeChapters) != 0 {
		response := gin.H{}
		response["error"] = "в томе ещё есть главы"

		for i := 0; i < len(volumeChapters); i++ {
			response[fmt.Sprintf("%d", i)] = volumeChapters[i]
		}

		c.AbortWithStatusJSON(409, response)

		return
	}

	if result := h.DB.Exec("DELETE FROM volumes WHERE id = ?", volumeID); result.Error != nil {
		log.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error})
		return
	}

	c.JSON(200, gin.H{"success": "том успешно удалён"})
}
