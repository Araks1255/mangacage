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
	DB                   *gorm.DB
	UsersProfilePictures *mongo.Collection
	NotificationsClient  pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, secretKey string, r *gin.Engine) {
	usersProfilePictures := client.Database("mangacage").Collection(mongodb.UsersProfilePicturesCollection)

	h := handler{
		DB:                   db,
		UsersProfilePictures: usersProfilePictures,
		NotificationsClient:  notificationsClient,
	}

	users := r.Group("/api/users")
	{
		users.GET("/", h.GetUsers)
		users.GET("/:id", h.GetUser)
		users.GET("/:id/profile-picture", h.GetUserProfilePicture)

		me := users.Group("/me")
		me.Use(middlewares.Auth(secretKey))
		{
			me.GET("/", h.GetMyProfile)
			me.GET("/profile-picture", h.GetMyProfilePicture)
			me.POST("/edited", h.EditProfile)
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, usersProfilePictures *mongo.Collection) handler {
	return handler{
		DB:                   db,
		NotificationsClient:  notificationsClient,
		UsersProfilePictures: usersProfilePictures,
	}
}
