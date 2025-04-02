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
		privateTeam.DELETE("/self/leave", h.LeaveTeam) // Глаголы в маршрутах вроде являются плохой практикой, но выдумывать какие-то "/team/members/self" с post запросом для "создания себя как участника команды" я не очень хочу.
		privateTeam.POST("/self/edited", h.EditTeam)   // Редактирование команды подразумевает создание отредактированной команды. Поэтому post. И ещё self добавил, для явного указания того, что команда ищется по юзеру, совершившему запрос (не уверен, что это нормальная практика, но вроде не плохая)

		privateTeam.POST("/:team/applications", h.SubmitTeamJoiningRequest) // Тут немного разная логика, поэтому объеденить в одну группу не выйдет
		privateTeam.DELETE("/:team/applications", h.CancelTeamJoiningRequest)
		privateTeam.GET("/self/applications", h.GetTeamJoiningApplications)
		privateTeam.GET("/applications/self", h.GetSelfJoiningApplications)
		privateTeam.POST("/self/applications/:candidate/accept", h.AcceptTeamJoiningApplication) // Тут я вообще не знаю как обойтись без глагола в маршруте, обновляется пользователь, а заявка удаляется. Но при этом логически, действие относится к заявке (она принимается)
		privateTeam.DELETE("/self/applications/:candidate", h.DeclineTeamJoiningApplication)
	}

	publicTeam := r.Group("api/teams")
	{
		publicTeam.GET("/:team/cover", h.GetTeamCover)
	}
}
