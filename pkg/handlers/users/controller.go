package users

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                               *gorm.DB
	UsersOnModerationProfilePictures *mongo.Collection
	UsersProfilePictures             *mongo.Collection
	NotificationsClient              pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, secretKey string, r *gin.Engine) {
	usersOnModerationProfilePictures := client.Database("mangacage").Collection(mongodb.UsersOnModerationProfilePicturesCollection)
	usersProfilePictures := client.Database("mangacage").Collection(mongodb.UsersProfilePicturesCollection)

	h := handler{
		DB:                               db,
		UsersOnModerationProfilePictures: usersOnModerationProfilePictures,
		UsersProfilePictures:             usersProfilePictures,
		NotificationsClient:              notificationsClient,
	}

	me := r.Group("/api/users/me")
	me.Use(middlewares.Auth(secretKey))
	{
		me.GET("/", h.GetMyProfile)
		me.GET("/profile-picture", h.GetMyProfilePicture)
		me.POST("/edited", h.EditProfile)
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
