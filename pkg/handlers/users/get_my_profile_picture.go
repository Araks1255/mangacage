package users

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetMyProfilePicture(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	filter := bson.M{"user_id": claims.ID}

	var result mongoModels.UserProfilePicture

	err := h.UsersProfilePictures.FindOne(c.Request.Context(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.AbortWithStatusJSON(404, gin.H{"error": "аватарка не найдена"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if len(result.ProfilePicture) == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при получении аватарки"})
		return
	}

	c.Data(200, "image/jpeg", result.ProfilePicture)
}
