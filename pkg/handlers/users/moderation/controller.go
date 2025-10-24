package moderation

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type handler struct{ DB *gorm.DB }

func RegisterRoutes(db *gorm.DB, secretKey string, r *gin.Engine) {
	h := handler{DB: db}

	moderation := r.Group("/api/users/me/moderation")
	moderation.Use(middlewares.Auth(secretKey))
	{
		profile := moderation.Group("/profile")
		{
			profile.DELETE("", h.CancelAppealForProfileChanges)
			profile.GET("/", h.GetMyProfileChangesOnModeration)
			profile.GET("/picture", h.GetMyProfilePictureOnModeration)
		}

		titles := moderation.Group("/titles")
		{
			titles.DELETE("/:id", h.CancelAppealForTitleModeration)
			titles.GET("/", h.GetMyTitlesOnModeration)
			titles.GET("/:id/cover", h.GetMyTitleOnModerationCover)
			titles.GET("/:id", h.GetMyTitleOnModeration)
		}

		chapters := moderation.Group("/chapters")
		{
			chapters.DELETE("/:id", h.CancelAppealForChapterModeration)
			chapters.GET("/", h.GetMyChaptersOnModeration)
			chapters.GET("/:id/page/:page", h.GetMyChapterOnModerationPage)
			chapters.GET("/:id", h.GetMyChapterOnModeration)
		}

		team := moderation.Group("/team")
		{
			team.DELETE("", h.CancelAppealForTeamModeration)
			team.GET("/", h.GetMyTeamOnModeration)
			team.GET("/cover", h.GetMyTeamOnModerationCover)
		}

		moderation.DELETE("/:entity/:id", h.CancelAppealForModeration)

		moderation.GET("/authors", h.GetMyAuthorsOnModeration)
		moderation.GET("/genres", h.GetMyGenresOnModeration)
		moderation.GET("/tags", h.GetMyTagsOnModeration)
	}
}

func NewHandler(db *gorm.DB) handler {
	return handler{DB: db}
}
