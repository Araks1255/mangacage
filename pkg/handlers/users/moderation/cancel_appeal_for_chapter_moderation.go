package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) CancelAppealForChapterModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredChapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "невалидный id главы на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		ChapterOnModerationID uint
		ChapterID             uint
	}

	if err = tx.Raw(
		`SELECT
			com.id AS chapter_on_moderation_id,
			c.id AS chapter_id
		FROM
			chapters_on_moderation AS com
			LEFT JOIN chapters AS c ON com.existing_id = c.id
		WHERE
			com.id = ? AND com.creator_id = ?`,
		desiredChapterOnModerationID, claims.ID,
	).Scan(&existing).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if existing.ChapterOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена среди ваших глав на модерации"})
		return
	}

	if result := tx.Exec("DELETE FROM chapters_on_moderation WHERE id = ?", existing.ChapterOnModerationID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if existing.ChapterID == 0 { // Страницы могут быть только у новой главы, у отредактированной (имеющией existing_id) - нет
		filter := bson.M{"chapter_on_moderation_id": existing.ChapterOnModerationID}

		if _, err = h.ChaptersPages.DeleteOne(c.Request.Context(), filter); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "обращение на модерацию главы успешно отменено"})
}
