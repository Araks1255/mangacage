package teams

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
	Collection *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	teamsCoversCollection := client.Database("mangacage").Collection("teams_covers")

	h := handler{
		DB: db,
		Collection: teamsCoversCollection,
	}

	team := r.Group("/teams")
	team.Use(middlewares.AuthMiddleware(secretKey))

	team.POST("/", h.CreateTeam)
	team.POST("/join/:team", h.JoinTeam)
	team.DELETE("/leave", h.LeaveTeam)
}
