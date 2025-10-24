package genres

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	NotificationsClient pb.SiteNotificationsClient
}

func RegisterRoutes(db *gorm.DB, notificationsClient pb.SiteNotificationsClient, secretKey string, r *gin.Engine) {
	h := handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}

	genres := r.Group("/api/genres")
	{
		genres.GET("/", middlewares.AuthOptional(secretKey), h.GetGenres)
		genres.POST("/", middlewares.Auth(secretKey), h.AddGenre)
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
