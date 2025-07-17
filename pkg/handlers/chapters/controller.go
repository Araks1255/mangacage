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
		chapters := api.Group("/chapters")
		{
			chapters.GET("/:id", h.GetChapter)
			chapters.GET("/", middlewares.AuthOptional(secretKey), h.GetChapters)
			chapters.GET("/:id/page/:page", h.GetChapterPage)

			chaptersAuth := chapters.Group("/")
			chaptersAuth.Use(middlewares.Auth(secretKey))
			{
				chaptersForTeamLeaders := chaptersAuth.Group("/")
				chaptersForTeamLeaders.Use(middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader", "translater"}))
				{
					chaptersForTeamLeaders.POST("/", h.CreateChapter)
					chaptersForTeamLeaders.POST("/:id/edited", h.EditChapter)
				}
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
