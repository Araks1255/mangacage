package moderation

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyProfilePictureOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userOnModerationID uint
	h.DB.Raw("SELECT id FROM users_on_moderation WHERE existing_id = ?", claims.ID).Scan(&userOnModerationID)
	if userOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет изменений профиля, ожидающих модерации"})
		return
	}

	var result struct {
		UserID         uint   `bson:"user_id"`
		ProfilePicture []byte `bson:"profile_picture"`
	}

	filter := bson.M{"user_on_moderation_id": userOnModerationID}

	if err := h.ProfilePictures.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет изменений аватарки на модерации"})
			return
		}
		log.Println(err)
		c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.ProfilePicture)
}
