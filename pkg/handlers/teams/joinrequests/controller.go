package joinrequests

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	SecretKey           string
	NotificationsCLient pb.SiteNotificationsClient
}

func RegisterRoutes(db *gorm.DB, secretKey string, notificationsCLient pb.SiteNotificationsClient, r *gin.Engine) {
	h := handler{
		DB:                  db,
		SecretKey:           secretKey,
		NotificationsCLient: notificationsCLient,
	}

	rolesRequire := []string{"team_leader", "vice_team_leader"}

	teamsAuth := r.Group("/api/teams")
	teamsAuth.Use(middlewares.Auth(secretKey))
	{
		teamsAuth.POST("/:id/join-requests", h.SubmitTeamJoinRequest)

		joinRequestsToMyTeam := teamsAuth.Group("/my/join-requests")
		{
			joinRequestsToMyTeam.GET("/", h.GetTeamJoinRequestsOfMyTeam)
			joinRequestsToMyTeam.GET("/:id", h.GetTeamJoinRequestOfMyTeam)

			joinRequestsForTeamLeaders := joinRequestsToMyTeam.Group("/")
			joinRequestsForTeamLeaders.Use(middlewares.RequireRoles(db, rolesRequire))
			{
				joinRequestsForTeamLeaders.POST("/:id/accept", h.AcceptTeamJoinRequest)
				joinRequestsForTeamLeaders.DELETE("/:id", h.DeclineTeamJoinRequest)
			}

		}

		joinRequests := teamsAuth.Group("/join-requests")
		{
			joinRequests.GET("/my", h.GetMyTeamJoinRequests)
			joinRequests.GET("/:id", h.GetMyTeamJoinRequest)
			joinRequests.DELETE("/:id", h.CancelTeamJoinRequest)
		}
	}
}

func NewHandler(db *gorm.DB, secretKey string, notificationsCLient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		SecretKey:           secretKey,
		NotificationsCLient: notificationsCLient,
	}
}
