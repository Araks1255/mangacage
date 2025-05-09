package volumes

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	NotificationsClient pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, notificationsClient pb.NotificationsClient, secretKey string, r *gin.Engine) {
	h := handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}

	privateVolume := r.Group("api/volumes/:id")
	privateVolume.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateVolume.POST("/edited", h.EditVolume)
		privateVolume.DELETE("/", h.DeleteVolume)
	}
	r.POST("/api/titles/:id/volumes", middlewares.AuthMiddleware(secretKey), h.CreateVolume)

	{
		r.GET("/api/volumes/:id", h.GetVolume)
		r.GET("/api/titles/:id/volumes", h.GetTitleVolumes)
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
