package moderation

import (
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyProfilePictureOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userOnModerationID sql.NullInt64

	if err := h.DB.Raw("SELECT id FROM users_on_moderation WHERE existing_id = ?", claims.ID).Scan(&userOnModerationID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !userOnModerationID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет изменений профиля, ожидающих модерации"})
		return
	}

	var result struct {
		UserID         uint   `bson:"user_id"`
		ProfilePicture []byte `bson:"profile_picture"`
	}

	filter := bson.M{"user_on_moderation_id": userOnModerationID.Int64}

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
