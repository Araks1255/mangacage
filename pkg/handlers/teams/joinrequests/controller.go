package joinrequests

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB        *gorm.DB
	SecretKey string
}

func RegisterRoutes(db *gorm.DB, secretKey string, r *gin.Engine) {
	h := handler{
		DB:        db,
		SecretKey: secretKey,
	}

	rolesRequire := []string{"team_leader", "ex_team_leader"}

	teamsAuth := r.Group("/api/teams")
	teamsAuth.Use(middlewares.Auth(secretKey))
	{
		teamsAuth.POST("/:id/join-requests", h.SubmitTeamJoinRequest)

		joinRequestsToMyTeam := teamsAuth.Group("/my/join-requests")
		{
			joinRequestsToMyTeam.GET("/", h.GetTeamJoinRequestsOfMyTeam)

			joinRequestsForTeamLeaders := joinRequestsToMyTeam.Group("/")
			joinRequestsForTeamLeaders.Use(middlewares.RequireRoles(db, rolesRequire))
			{
				joinRequestsForTeamLeaders.POST("/:id/accept", h.AcceptTeamJoinRequest)
				joinRequestsForTeamLeaders.DELETE("/:id", h.DeclineTeamJoinRequest)
			}

		}

		joinRequests := teamsAuth.Group("/join-requests")
		{
			joinRequests.GET("my", h.GetMyTeamJoinRequests)
			joinRequests.DELETE("/:id", h.CancelTeamJoinRequest)
		}
	}
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		DB: db,
	}
}
