package users

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB         *gorm.DB
	Collection *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	usersOnModerationProfilePictures := client.Database("mangacage").Collection("users_on_moderation_profile_pictures")

	h := handler{
		DB:         db,
		Collection: usersOnModerationProfilePictures,
	}

	privateUser := r.Group("/home")
	privateUser.Use(middlewares.AuthMiddleware(secretKey))

	privateUser.GET("/viewed_chapters/inf", h.GetViewedChapters)
	privateUser.GET("profile/inf", h.GetSelfProfile)
	privateUser.GET("/profile/profile_picture", h.GetSelfProfilePicture)
	privateUser.POST("/profile/edit", h.EditProfile)
}
