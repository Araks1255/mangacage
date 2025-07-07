package titles

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	TitlesCovers        *mongo.Collection
	NotificationsClient pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, secretKey string, r *gin.Engine) {
	mongoDB := client.Database("mangacage")

	titlesCoversCollection := mongoDB.Collection(mongodb.TitlesCoversCollection)

	h := handler{
		DB:                  db,
		TitlesCovers:        titlesCoversCollection,
		NotificationsClient: notificationsClient,
	}

	titles := r.Group("/api/titles")
	{
		titles.GET("/:id", middlewares.AuthOptional(secretKey), h.GetTitle)
		titles.GET("/", middlewares.AuthOptional(secretKey), h.GetTitles)
		titles.GET("/:id/cover", h.GetTitleCover)

		titlesAuth := titles.Group("/")
		titlesAuth.Use(middlewares.Auth(secretKey))
		{
			titlesAuth.POST("/:id/subscriptions", h.SubscribeToTitle)
			titlesAuth.POST("/", h.CreateTitle)

			rates := titlesAuth.Group("/:id/rate")
			{
				rates.POST("/", h.RateTitle)
				rates.DELETE("/", h.DeleteTitleRate)
			}

			titlesForTeamLeaders := titlesAuth.Group("/:id")
			titlesForTeamLeaders.Use(middlewares.RequireRoles(db, []string{"team_leader"}))
			{
				titlesForTeamLeaders.PATCH("/translate", h.TranslateTitle)
				titlesForTeamLeaders.DELETE("/quit-translating", h.QuitTranslatingTitle)
				// titlesForTeamLeaders.DELETE("/", h.DeleteTitle)
			}

			titlesForExTeamLeaders := titlesAuth.Group("/:id")
			titlesForExTeamLeaders.Use(middlewares.RequireRoles(db, []string{"team_leader", "ex_team_leader"}))
			{
				titlesForExTeamLeaders.POST("/edited", h.EditTitle)
			}

			translateRequests := titlesAuth.Group("/translate-requests")
			{
				translateRequests.GET("/", h.GetTitleTranslateRequests)
				translateRequests.DELETE(
					"/:id",
					middlewares.RequireRoles(db, []string{"team_leader"}),
					h.CancelTitleTranslateRequest,
				)
			}
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, titlesCovers *mongo.Collection) handler {
	return handler{
		DB:                  db,
		TitlesCovers:        titlesCovers,
		NotificationsClient: notificationsClient,
	}
}
