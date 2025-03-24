package moderation

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetSelfProfilePictureOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var result struct {
		UserID         uint   `bson:"user_id"`
		ProfilePicture []byte `bson:"profile_picture"`
	}

	filter := bson.M{"user_id": claims.ID}

	if err := h.ProfilePictures.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.ProfilePicture)
}
