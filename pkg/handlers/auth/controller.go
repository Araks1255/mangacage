package auth

import (
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	NotificationsClient pb.SiteNotificationsClient
	SecretKey           string
	Host                string
}

func RegisterRoutes(db *gorm.DB, notificationsClient pb.SiteNotificationsClient, secretKey, host string, r *gin.Engine) {
	h := handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
		SecretKey:           secretKey,
		Host:                host,
	}

	api := r.Group("/api")
	{
		api.POST("/signup", h.Signup)
		api.POST("/login", h.Login)
		api.POST("/logout", h.Logout)
	}
}
