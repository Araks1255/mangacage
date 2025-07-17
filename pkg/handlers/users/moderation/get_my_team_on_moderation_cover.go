package moderation

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyTeamOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var result mongoModels.TeamOnModerationCover

	filter := bson.M{"creator_id": claims.ID}

	if err := h.TeamsCovers.FindOne(c.Request.Context(), filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.AbortWithStatusJSON(404, gin.H{"error": "обложка команды на модерации не найдена"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
