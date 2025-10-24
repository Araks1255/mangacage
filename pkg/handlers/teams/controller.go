package teams

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	PathToMediaDir string
	NotificationsClient pb.SiteNotificationsClient
}

func RegisterRoutes(db *gorm.DB, pathToMediaDir string, notificationsClient pb.SiteNotificationsClient, secretKey string, r *gin.Engine) {
	h := handler{
		DB:                  db,
		PathToMediaDir: pathToMediaDir,
		NotificationsClient: notificationsClient,
	}

	rolesRequired := []string{"team_leader"}

	teams := r.Group("/api/teams")
	{
		teams.GET("/:id/cover", h.GetTeamCover)
		teams.GET("/:id/", h.GetTeam)
		teams.GET("/", h.GetTeams)

		teamsAuth := teams.Group("/")
		teamsAuth.Use(middlewares.Auth(secretKey))
		{
			teamsAuth.POST("/", h.CreateTeam)

			teamsForTeamLeaders := teamsAuth.Group("/my")
			teamsForTeamLeaders.Use(middlewares.RequireRoles(db, rolesRequired))
			{
				teamsForTeamLeaders.POST("/edited", h.EditTeam) // Тут создание отредактированной команды, поэтому post
				teamsForTeamLeaders.DELETE("/", h.DeleteTeam)
			}
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
