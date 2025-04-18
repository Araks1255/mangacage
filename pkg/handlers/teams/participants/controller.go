package participants

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

	privateParticipants := r.Group("/api/teams/my/participants")
	privateParticipants.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateParticipants.PATCH("/:id/role", h.ChangeParticipantRole)
		privateParticipants.DELETE("/me", h.LeaveTeam)
	}

	publicParticipants := r.Group("/api/teams/:id/participants")
	{
		publicParticipants.GET("/", h.GetTeamParticipants)
	}
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		DB: db,
	}
}
