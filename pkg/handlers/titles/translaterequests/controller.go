package translaterequests

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

func RegisterRoutes(db *gorm.DB, secretKey string, notificationsClient pb.SiteNotificationsClient, r *gin.Engine) {
	h := handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}

	titles := r.Group("/api/titles")
	titles.Use(middlewares.Auth(secretKey), middlewares.RequireRoles(db, []string{"team_leader"}))
	{
		titles.PATCH("/:id/translate", h.TranslateTitle)
		titles.DELETE("/:id/quit-translating", h.QuitTranslatingTitle)
	}

	translateRequests := r.Group("/api/users/me/titles-translate-requests")
	translateRequests.Use(middlewares.Auth(secretKey))
	{
		translateRequests.GET("/", h.GetMyTitleTranslateRequests)
		translateRequests.GET("/:id", h.GetMyTitleTranslateRequest)
		translateRequests.DELETE("/:id", middlewares.RequireRoles(db, []string{"team_leader"}), h.CancelTitleTranslateRequest)
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
