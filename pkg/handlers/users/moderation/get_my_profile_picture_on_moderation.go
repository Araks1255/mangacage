package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyProfilePictureOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var result mongoModels.UserOnModerationProfilePicture

	filter := bson.M{"creator_id": claims.ID}

	if err := h.ProfilePictures.FindOne(c.Request.Context(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет изменений аватарки на модерации"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.ProfilePicture)
}
