package auth

import (
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	NotificationsClient pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, r *gin.Engine) {
	h := handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}

	api := r.Group("/api")
	{
		api.POST("/signup", h.Signup)
		api.POST("/login", h.Login)
		api.POST("/logout", h.Logout)
	}
}
