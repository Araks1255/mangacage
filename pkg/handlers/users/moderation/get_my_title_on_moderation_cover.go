package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyTitleOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	var doesTitleOnModerationExist bool

	if err := h.DB.Raw(
		"SELECT EXISTS(SELECT 1 FROM titles_on_moderation WHERE id = ? AND creator_id = ?)",
		titleOnModerationID, claims.ID,
	).Scan(&doesTitleOnModerationExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesTitleOnModerationExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди ваших тайтлов на модерации"})
		return
	}

	var result struct {
		TitleOnModerationID uint   `bson:"title_on_moderation"`
		Cover               []byte `bson:"cover"`
	}

	filter := bson.M{"title_on_moderation_id": titleOnModerationID}

	if err := h.TitlesCovers.FindOne(c.Request.Context(), filter).Decode(&result); err != nil {
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
