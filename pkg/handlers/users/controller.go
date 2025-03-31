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

	privateUser := r.Group("/api/home")
	privateUser.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateUser.GET("/viewed_chapters", h.GetViewedChapters)

		profile := privateUser.Group("/profile")
		{
			profile.GET("/", h.GetSelfProfile)
			profile.GET("/picture", h.GetSelfProfilePicture)
			profile.POST("/edited", h.EditProfile)
		}
	}
}
