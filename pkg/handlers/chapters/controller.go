package chapters

import (
	cpc "github.com/Araks1255/mangacage/internal/workers/chapters_pages_compressor"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type handler struct {
	DB                      *gorm.DB
	PathToMediaDir          string
	ChaptersPagesCompressor *cpc.ChaptersPagesCompressor
	NotificationsClient     pb.SiteNotificationsClient
}

func RegisterRoutes(db *gorm.DB, pathToMediaDir string, compressor *cpc.ChaptersPagesCompressor, notificationsClient pb.SiteNotificationsClient, secretKey string, r *gin.Engine) {
	h := handler{
		DB:                      db,
		PathToMediaDir:          pathToMediaDir,
		ChaptersPagesCompressor: compressor,
		NotificationsClient:     notificationsClient,
	}

	api := r.Group("/api")
	{
		chapters := api.Group("/chapters")
		{
			chapters.GET("/:id", h.GetChapter)
			chapters.GET("", middlewares.AuthOptional(secretKey), h.GetChapters)
			chapters.GET("/:id/page/:page", h.GetChapterPage)

			chaptersAuth := chapters.Group("")
			chaptersAuth.Use(middlewares.Auth(secretKey))
			{
				chaptersForTeamLeaders := chaptersAuth.Group("")
				chaptersForTeamLeaders.Use(middlewares.RequireRoles(db, []string{"team_leader", "vice_team_leader", "translater"}))
				{
					chaptersForTeamLeaders.POST("", h.CreateChapter)
					chaptersForTeamLeaders.POST("/:id/edited", h.EditChapter)
				}
			}
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
