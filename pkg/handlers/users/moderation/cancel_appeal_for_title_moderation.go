package moderation

import (
	"context"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) CancelAppealForTitleModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingTitleOnModerationID uint
	tx.Raw("SELECT id FROM titles_on_moderation WHERE id = ? AND creator_id = ?", desiredTitleOnModerationID, claims.ID).Scan(&existingTitleOnModerationID)
	if existingTitleOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден в ваших тайтлах на модерации"})
		return
	}

	if result := tx.Exec("DELETE FROM titles_on_moderation WHERE id = ?", existingTitleOnModerationID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	filter := bson.M{"title_on_moderation_id": existingTitleOnModerationID}

	if _, err := h.TitlesCovers.DeleteOne(context.TODO(), filter); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваше обращение на модерацию тайтла успешно отменено"})
}
