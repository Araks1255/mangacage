package teams

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB          *gorm.DB
	TeamsCovers *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, secretKey string, r *gin.Engine) {
	teamsCoversCollection := client.Database("mangacage").Collection("teams_covers")

	h := handler{
		DB:          db,
		TeamsCovers: teamsCoversCollection,
	}

	rolesRequired := []string{"team_leader"}

	teams := r.Group("/api/teams")
	{
		teams.GET("/:id/cover", h.GetTeamCover)
		teams.GET("/:id/", h.GetTeam)

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

func NewHandler(db *gorm.DB, teamsCovers *mongo.Collection) handler {
	return handler{
		DB:          db,
		TeamsCovers: teamsCovers,
	}
}
