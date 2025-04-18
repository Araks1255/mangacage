package joinrequests

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	h := handler{DB: db}

	teamJoinRequests := r.Group("/api/teams/my/join-requests") // Для тех пользователей, кто уже в команде (просмотр текущих заявок в свою команду, одобрение/отклонение)
	teamJoinRequests.Use(middlewares.AuthMiddleware(secretKey))
	{
		teamJoinRequests.POST("/:id/accept", h.AcceptTeamJoinRequest)
		teamJoinRequests.DELETE("/:id", h.DeclineTeamJoinRequest)
		teamJoinRequests.GET("/my_team", h.GetTeamJoinRequestsOfMyTeam)
	}

	usersTeamJoinRequests := r.Group("/api/teams/join-requests") // Для пользователей, не находящихся в команде (отмена и получение своих заявок)
	usersTeamJoinRequests.Use(middlewares.AuthMiddleware(secretKey))
	{
		usersTeamJoinRequests.GET("/my", h.GetMyTeamJoinRequests)
		usersTeamJoinRequests.DELETE("/:id", h.CancelTeamJoinRequest)
	}

	r.POST("/api/teams/:id/join-requests", middlewares.AuthMiddleware(secretKey), h.SubmitTeamJoinRequest) // Это подача заявки. Она по логике ко второй группе относится, но у неё в маршруте id конкретной команды, который по-моему не совсем уместно было бы выносить в тело запроса
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		DB: db,
	}
}
