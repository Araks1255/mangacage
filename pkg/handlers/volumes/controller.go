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

	api := r.Group("/api")
	{
		titles := api.Group("/titles/:id")
		{
			titles.GET("/volumes", h.GetTitleVolumes)

			titles.POST(
				"/volumes",
				middlewares.Auth(secretKey),
				middlewares.RequireRoles(db, []string{"team_leader"}),
				h.CreateVolume,
			)
		}

		volumes := api.Group("/volumes")
		{
			volumes.GET("/:id", h.GetVolume)

			volumesAuth := volumes.Group("/")
			volumesAuth.Use(middlewares.Auth(secretKey))
			{
				// volumesAuth.DELETE(
				// 	"/:id",
				// 	middlewares.RequireRoles(db, []string{"team_leader"}),
				// 	h.DeleteVolume,
				// )

				volumesAuth.POST(
					"/:id/edited",
					middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader"}),
					h.EditVolume,
				)
			}
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
