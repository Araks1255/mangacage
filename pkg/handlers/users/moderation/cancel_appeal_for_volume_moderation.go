package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelAppealForVolumeModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")
	volume := c.Param("volume")

	tx := h.DB.Begin()
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}

	var volumeID uint
	tx.Raw(
		`SELECT v.id FROM volumes_on_moderation AS v
		INNER JOIN titles ON titles.id = v.title_id
		WHERE v.creator_id = ?
		AND titles.name = ?
		AND v.name = ?`,
		claims.ID, title, volume,
	).Scan(&volumeID)

	if volumeID == 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	if result := tx.Exec("DELETE FROM volumes_on_moderation WHERE id = ?", volumeID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "обращение на модерацию тома успешно отменено"})
}
