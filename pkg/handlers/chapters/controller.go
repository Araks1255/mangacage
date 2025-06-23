package chapters

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	ChaptersPages       *mongo.Collection
	NotificationsClient pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, secretKey string, r *gin.Engine) {
	chapterPagesCollection := client.Database("mangacage").Collection("chapters_pages")

	h := handler{
		DB:                  db,
		ChaptersPages:       chapterPagesCollection,
		NotificationsClient: notificationsClient,
	}

	api := r.Group("/api")
	{
		api.POST(
			"/volumes/:id/chapters",
			middlewares.Auth(secretKey),
			middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader"}),
			h.CreateChapter,
		)

		chapters := api.Group("/chapters")
		{
			chapters.GET("/:id", h.GetChapter)
			chapters.GET("/", middlewares.AuthOptional(secretKey), h.GetChapters)
			chapters.GET("/:id/page/:page", h.GetChapterPage)

			chaptersAuth := chapters.Group("/")
			chaptersAuth.Use(middlewares.Auth(secretKey))
			{
				// chaptersAuth.DELETE(
				// 	"/",
				// 	middlewares.RequireRoles(db, []string{"team_leader"}),
				// 	h.DeleteChapter,
				// )

				chaptersAuth.POST( // Тут идёт создание отредактированной главы (прямо отдельная сущность в отдельной таблице базы данных), поэтому post а не put
					"/",
					middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader"}),
					h.EditChapter,
				)
			}
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, chaptersPages *mongo.Collection) handler {
	return handler{
		DB:                  db,
		ChaptersPages:       chaptersPages,
		NotificationsClient: notificationsClient,
	}
}
