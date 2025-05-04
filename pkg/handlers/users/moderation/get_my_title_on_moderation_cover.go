package moderation

import (
	"context"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyTitleOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	var existingTitleOnModerationID uint
	h.DB.Raw("SELECT id FROM titles_on_moderation WHERE id = ? AND creator_id = ?", desiredTitleOnModerationID, claims.ID).Scan(&existingTitleOnModerationID)
	if existingTitleOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди ваших тайтлов на модерации"})
		return
	}

	var result struct {
		TitleOnModerationID uint   `bson:"title_on_moderation"`
		Cover               []byte `bson:"cover"`
	}

	filter := bson.M{"title_on_moderation_id": existingTitleOnModerationID}

	if err := h.TitlesCovers.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(404, gin.H{"error": "обложка тайтла на модерации не найдена"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Cover) == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при получении обложки тайтла на модерации"})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
