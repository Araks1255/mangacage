package users

import (
	"github.com/Araks1255/mangacage/pkg/constants"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                               *gorm.DB
	UsersOnModerationProfilePictures *mongo.Collection
	UsersProfilePictures             *mongo.Collection
	NotificationsClient              pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	usersOnModerationProfilePictures := client.Database("mangacage").Collection(constants.UsersOnModerationProfilePicturesCollection)
	usersProfilePictures := client.Database("mangacage").Collection(constants.UsersProfilePicturesCollection)

	h := handler{
		DB:                               db,
		UsersOnModerationProfilePictures: usersOnModerationProfilePictures,
		UsersProfilePictures:             usersProfilePictures,
		NotificationsClient:              notificationsClient,
	}

	privateUser := r.Group("/api/users/me")
	privateUser.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateUser.GET("/", h.GetMyProfile)
		privateUser.GET("/profile-picture", h.GetMyProfilePicture)
		privateUser.POST("/edited", h.EditProfile)
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, usersProfilePictures, usersOnModerationProfilePictures *mongo.Collection) handler {
	return handler{
		DB:                               db,
		NotificationsClient:              notificationsClient,
		UsersProfilePictures:             usersProfilePictures,
		UsersOnModerationProfilePictures: usersOnModerationProfilePictures,
	}
}
