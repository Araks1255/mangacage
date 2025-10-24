package titles

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	PathToMediaDir      string
	NotificationsClient pb.SiteNotificationsClient
}

func RegisterRoutes(db *gorm.DB, pathToMediaDir string, notificationsClient pb.SiteNotificationsClient, secretKey string, r *gin.Engine) {
	h := handler{
		DB:                  db,
		PathToMediaDir:      pathToMediaDir,
		NotificationsClient: notificationsClient,
	}

	titles := r.Group("/api/titles")
	{
		titles.GET("/:id", middlewares.AuthOptional(secretKey), h.GetTitle)
		titles.GET("/", middlewares.AuthOptional(secretKey), h.GetTitles)
		titles.GET("/:id/cover", h.GetTitleCover)

		titlesAuth := titles.Group("/")
		titlesAuth.Use(middlewares.Auth(secretKey))
		{
			titlesAuth.POST("/:id/subscribe", h.SubscribeToTitle)
			titlesAuth.DELETE("/:id/unsubscribe", h.UnSubscribeFromTitle)
			titlesAuth.POST("/", h.CreateTitle)

			rates := titlesAuth.Group("/:id/rate")
			{
				rates.POST("/", h.RateTitle)
				rates.DELETE("/", h.DeleteTitleRate)
			}

			titlesAuth.POST("/:id/edited", middlewares.RequireRoles(db, []string{"team_leader"}), h.EditTitle)
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
