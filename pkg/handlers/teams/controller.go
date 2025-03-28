package teams

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                      *gorm.DB
	TeamsOnModerationCovers *mongo.Collection
	TeamsCovers *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	teamsOnModerationCoversCollection := client.Database("mangacage").Collection("teams_on_moderation_covers")
	teamsCoversCollection := client.Database("mangacage").Collection("teams_covers")

	h := handler{
		DB:                      db,
		TeamsOnModerationCovers: teamsOnModerationCoversCollection,
		TeamsCovers: teamsCoversCollection,
	}

	privateTeam := r.Group("/teams")
	privateTeam.Use(middlewares.AuthMiddleware(secretKey))

	privateTeam.POST("/", h.CreateTeam)
	privateTeam.POST("/join/:team", h.JoinTeam)
	privateTeam.DELETE("/leave", h.LeaveTeam)
	privateTeam.POST("/edit", h.EditTeam)

	publicTeam := r.Group("/teams")

	publicTeam.GET("/:team/cover", h.GetTeamCover)
}
