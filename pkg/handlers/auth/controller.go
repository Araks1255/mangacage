package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB         *gorm.DB
	Collection *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	usersProfilePicturesCollection := client.Database("mangacage").Collection("users_on_moderation_profile_pictures")

	h := handler{
		DB:         db,
		Collection: usersProfilePicturesCollection,
	}

	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
}
