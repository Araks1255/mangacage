package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelAppealForChapterModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")
	volume := c.Param("volume")
	chapter := c.Param("chapter")

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var chapterID uint
	tx.Raw(
		`SELECT c.id FROM chapters_on_moderation AS c
		INNER JOIN volumes ON volumes.id = c.volume_id
		INNER JOIN titles ON titles.id = volumes.title_id
		WHERE titles.name = ?
		AND volumes.name = ?
		AND c.name = ?
		AND c.creator_id = ?`,
		title, volume, chapter, claims.ID,
	).Scan(&chapterID)

	if chapterID == 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	if result := tx.Exec("DELETE FROM chapters_on_moderation WHERE id = ?", chapterID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "обращение на модерацию успешно отменено"})
}
