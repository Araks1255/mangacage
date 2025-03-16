package users

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetSelfProfilePicture(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	filter := bson.M{"user_id": claims.ID}

	var result struct {
		UserID         uint   `bson:"user_id"`
		ProfilePicture []byte `bson:"profile_picture"`
	}

	err := h.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.ProfilePicture) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "аватарка не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.ProfilePicture)
}
