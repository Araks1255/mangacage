package authors

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

	authors := r.Group("/api/authors")
	{
		authors.GET("/:id", h.GetAuthor)
		authors.GET("/", h.GetAuthors)

		authorsAuth := authors.Group("/")
		authorsAuth.Use(middlewares.Auth(secretKey))
		{
			authorsAuth.POST("/", h.AddAuthor)
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
