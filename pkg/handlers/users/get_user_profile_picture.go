package users

import (
	"errors"
	"strconv"

	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetUserProfilePicture(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id пользователя"})
		return
	}

	var res mongoModels.UserProfilePicture

	filter := bson.M{"user_id": userID, "visible": true}

	if err := h.UsersProfilePictures.FindOne(c.Request.Context(), filter).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.AbortWithStatusJSON(404, gin.H{"error": "аватарка не найдена"})
			return
		}
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", res.ProfilePicture)
}
