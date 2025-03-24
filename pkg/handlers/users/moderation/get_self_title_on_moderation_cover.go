package moderation

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetSelfTitleOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")

	var titleID uint
	h.DB.Raw("SELECT id FROM titles_on_moderation WHERE name = ? AND creator_id = ?", title, claims.ID).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден в ваших тайтлах на модерации"})
		return
	}

	var titleCover struct {
		TitleID uint   `bson:"title_id"`
		Cover   []byte `bson:"cover"`
	}

	filter := bson.M{"title_id": titleID}

	if err := h.TitlesCovers.FindOne(context.TODO(), filter).Decode(&titleCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", titleCover.Cover)
}
