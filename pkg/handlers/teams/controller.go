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
	TeamsCovers             *mongo.Collection
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
		TeamsCovers:             teamsCoversCollection,
	}

	privateTeam := r.Group("api/teams")
	privateTeam.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateTeam.POST("/", h.CreateTeam)
		privateTeam.PUT("/:team/join", h.JoinTeam) // Глаголы в маршрутах вроде являются плохой практикой, но выдумывать какие-то "/team/members/self" с post запросом для "создания себя как участника команды" я не очень хочу.
		privateTeam.DELETE("/self/leave", h.LeaveTeam)
		privateTeam.POST("/self/edited", h.EditTeam) // Редактирование команды подразумевает создание отредактированной команды. Поэтому post. И ещё self добавил, для явного указания того, что команда ищется по юзеру, совершившему запрос (не уверен, что это нормальная практика, но вроде не плохая)
	}

	publicTeam := r.Group("api/teams")
	{
		publicTeam.GET("/:team/cover", h.GetTeamCover)
	}
}
