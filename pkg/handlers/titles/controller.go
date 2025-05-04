package titles

import (
	"github.com/Araks1255/mangacage/pkg/constants"
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

	titlesCoversCollection := mongoDB.Collection(constants.TitlesCoversCollection)
	titlesOnModerationCovers := mongoDB.Collection(constants.TitlesOnModerationCoversCollection)

	h := handler{
		DB:                       db,
		TitlesCovers:             titlesCoversCollection,
		TitlesOnModerationCovers: titlesOnModerationCovers,
		NotificationsClient:      notificationsClient,
	}

	privateTitle := r.Group("/api/titles")
	privateTitle.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateTitle.POST("/", h.CreateTitle)
		privateTitle.PATCH("/:id/translate", h.TranslateTitle)
		privateTitle.POST("/:id/subscriptions", h.SubscribeToTitle)
		privateTitle.POST("/:id/edited", h.EditTitle)
		privateTitle.DELETE("/:id", h.DeleteTitle)
		privateTitle.PATCH("/:id/quit-translating", h.QuitTranslatingTitle)
	}

	publicTitle := r.Group("/api/titles")
	{
		publicTitle.GET("/:id/cover", h.GetTitleCover)
		publicTitle.GET("/most-popular", h.GetMostPopularTitles)
		publicTitle.GET("/recently-updated", h.GetRecentlyUpdatedTitles)
		publicTitle.GET("/new", h.GetNewTitles)
		publicTitle.GET("/:id", h.GetTitle)
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
