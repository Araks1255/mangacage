package chapters

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                        *gorm.DB
	ChaptersOnModerationPages *mongo.Collection
	ChaptersPages             *mongo.Collection
	NotificationsClient       pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, secretKey string, r *gin.Engine) {
	chaptersOnModerationPagesCollection := client.Database("mangacage").Collection("chapters_on_moderation_pages")
	chapterPagesCollection := client.Database("mangacage").Collection("chapters_pages")

	h := handler{
		DB:                        db,
		ChaptersOnModerationPages: chaptersOnModerationPagesCollection,
		ChaptersPages:             chapterPagesCollection,
		NotificationsClient:       notificationsClient,
	}

	api := r.Group("/api")
	{
		volumes := api.Group("/volumes/:id")
		{
			volumes.GET("/chapters", h.GetVolumeChapters)
			volumes.POST(
				"/chapters",
				middlewares.Auth(secretKey),
				middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader"}),
				h.CreateChapter,
			)
		}

		chapters := api.Group("/chapters/:id")
		{
			chapters.GET("/", h.GetChapter)
			chapters.GET("/page/:page", h.GetChapterPage)

			chaptersAuth := chapters.Group("/")
			chaptersAuth.Use(middlewares.Auth(secretKey))
			{
				chaptersAuth.DELETE(
					"/",
					middlewares.RequireRoles(db, []string{"team_leader"}),
					h.DeleteChapter,
				)

				chaptersAuth.POST( // Тут идёт создание отредактированной главы (прямо отдельная сущность в отдельной таблице базы данных), поэтому post а не put
					"/",
					middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader"}),
					h.EditChapter,
				)
			}
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, chaptersOnModerationPages, chaptersPages *mongo.Collection) handler {
	return handler{
		DB:                        db,
		ChaptersOnModerationPages: chaptersOnModerationPages,
		ChaptersPages:             chaptersPages,
		NotificationsClient:       notificationsClient,
	}
}
