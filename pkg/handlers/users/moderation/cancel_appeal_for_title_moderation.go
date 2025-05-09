package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) CancelAppealForTitleModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesTitleOnModerationExist bool

	if err := tx.Raw(
		"SELECT EXISTS(SELECT 1 FROM titles_on_moderation WHERE id = ? AND creator_id = ?)",
		titleOnModerationID, claims.ID,
	).Scan(&doesTitleOnModerationExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error":err.Error()})
		return
	}

	if !doesTitleOnModerationExist {
		c.AbortWithStatusJSON(404, gin.H{"error":"тайтл не найден среди ваших тайтлов на модерации"})
		return
	}

	if result := tx.Exec("DELETE FROM titles_on_moderation WHERE id = ?", titleOnModerationID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	filter := bson.M{"title_on_moderation_id": titleOnModerationID}

	if _, err := h.TitlesCovers.DeleteOne(c.Request.Context(), filter); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваше обращение на модерацию тайтла успешно отменено"})
}
