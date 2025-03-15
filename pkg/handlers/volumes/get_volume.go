package volumes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) GetVolume(c *gin.Context) {
	title := c.Param("title")
	desiredVolume := c.Param("volume")

	var desiredVolumeID uint
	h.DB.Raw(`SELECT volumes.id FROM volumes
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE lower(titles.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND NOT volumes.on_moderation`,
		title, desiredVolume).
		Scan(&desiredVolumeID)

	if desiredVolumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	var volume struct {
		gorm.Model
		Name        string
		Description string
		Title       string
		Creator     string
		Team        string
	}

	h.DB.Raw(`SELECT volumes.id, volumes.created_at, volumes.updated_at, volumes.deleted_at, volumes.name, volumes.description,
		titles.name AS title, users.user_name AS creator, teams.name AS team FROM volumes
		INNER JOIN titles ON volumes.title_id = titles.id
		INNER JOIN users ON volumes.creator_id = users.id
		INNER JOIN teams ON titles.team_id = teams.id
		WHERE volumes.id = ?`, desiredVolumeID).Scan(&volume)

	c.JSON(200, &volume)
}
