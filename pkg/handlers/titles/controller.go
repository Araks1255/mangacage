package titles

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                       *gorm.DB
	TitlesCovers             *mongo.Collection
	TitlesOnModerationCovers *mongo.Collection
	NotificationsClient      pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	mongoDB := client.Database("mangacage")

	titlesCoversCollection := mongoDB.Collection(mongodb.TitlesCoversCollection)
	titlesOnModerationCovers := mongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

	h := handler{
		DB:                       db,
		TitlesCovers:             titlesCoversCollection,
		TitlesOnModerationCovers: titlesOnModerationCovers,
		NotificationsClient:      notificationsClient,
	}

	titles := r.Group("/api/titles")
	{
		titles.GET("/:id/cover", h.GetTitleCover)
		titles.GET("/most-popular", h.GetMostPopularTitles)
		titles.GET("/new", h.GetNewTitles)
		titles.GET("/:id", h.GetTitle)

		titlesAuth := titles.Group("/")
		titlesAuth.Use(middlewares.Auth(secretKey))
		{
			titlesAuth.POST("/:id/subscriptions", h.SubscribeToTitle)
			titlesAuth.POST("/", h.CreateTitle)

			titlesForTeamLeaders := titlesAuth.Group("/:id")
			titlesForTeamLeaders.Use(middlewares.RequireRoles(db, []string{"team_leader"}))
			{
				titlesForTeamLeaders.PATCH("/translate", h.TranslateTitle)
				titlesForTeamLeaders.DELETE("/quit-translating", h.QuitTranslatingTitle)
				titlesForTeamLeaders.DELETE("/", h.DeleteTitle)
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

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, titlesCovers, titlesOnModerationCovers *mongo.Collection) handler {
	return handler{
		DB:                       db,
		TitlesCovers:             titlesCovers,
		TitlesOnModerationCovers: titlesOnModerationCovers,
		NotificationsClient:      notificationsClient,
	}
}
