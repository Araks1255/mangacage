package participants

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, secretKey string, r *gin.Engine) {
	h := handler{DB: db}

	teams := r.Group("/api/teams")
	{
		participantsOfMyTeam := teams.Group("/my/participants")
		participantsOfMyTeam.Use(middlewares.Auth(secretKey))
		{
			participantsOfMyTeam.DELETE("/me", h.LeaveTeam)                // Это покидание команды, роли никакие не нужны
			participantsOfMyTeam.POST(":id/roles", h.AddRoleToParticipant) // Это 2 хэндлера управления ролями участников. Они требуют более гибкой работы над ролями пользователя, так что middleware в их случае не используется
			participantsOfMyTeam.DELETE("/:id/roles", h.DeleteParticipantRole)
			participantsOfMyTeam.DELETE("/:id", middlewares.RequireRoles(db, []string{"team_leader"}), h.ExcludeParticipant) // Исключать может только лидер
		}
	}
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		DB: db,
	}
}
