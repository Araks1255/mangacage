package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
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

	var result mongoModels.TitleOnModerationCover

	filter := bson.M{"title_on_moderation_id": titleOnModerationID, "creator_id": claims.ID}

	if err := h.TitlesCovers.FindOne(c.Request.Context(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(404, gin.H{"error": "обложка тайтла на модерации не найдена"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
