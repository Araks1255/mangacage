package moderation

import (
	"context"
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	defer tx.Rollback()

	var (
		chapterID             sql.NullInt64
		chapterOnModerationID uint
	)

	row := tx.Raw(
		`SELECT c.existing_id, c.id FROM chapters_on_moderation AS c
		INNER JOIN volumes AS v ON v.id = c.volume_id
		INNER JOIN titles AS t ON t.id = v.title_id
		WHERE t.name = ? AND v.name = ? AND c.name = ? AND c.creator_id = ?`,
		title, volume, chapter, claims.ID,
	).Row()

	if err := row.Scan(&chapterID, &chapterOnModerationID); err != nil {
		log.Println(err)
	}

	if chapterOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена в ваших обращениях на модерацию"})
		return
	}

	if result := tx.Exec("DELETE FROM chapters_on_moderation WHERE id = ?", chapterID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	var filter bson.M

	if chapterID.Valid {
		filter = bson.M{"chapter_id": chapterID}
	} else {
		filter = bson.M{"chapter_on_moderation_id": chapterOnModerationID}
	}

	if _, err := h.ChaptersPages.DeleteOne(context.TODO(), filter); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "обращение на модерацию успешно отменено"})
}
