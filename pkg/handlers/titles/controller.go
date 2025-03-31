package titles

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB         *gorm.DB
	Collection *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	titlesCoversCollection := client.Database("mangacage").Collection("titles_on_moderation_covers")

	h := handler{
		DB:         db,
		Collection: titlesCoversCollection,
	}

	privateTitle := r.Group("/api/titles")
	privateTitle.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateTitle.POST("/", h.CreateTitle)
		privateTitle.PUT("/:title/translate", h.TranslateTitle)
		privateTitle.POST("/:title/subscribe", h.SubscribeToTitle)
		privateTitle.POST("/:title/edited", h.EditTitle)
		privateTitle.DELETE("/:title", h.DeleteTitle)
		privateTitle.PUT("/:title/quit", h.QuitTranslatingTitle)
	}

	publicTitle := r.Group("/api/titles")
	{
		publicTitle.GET("/:title/cover", h.GetTitleCover)
		publicTitle.GET("/most_popular", h.GetMostPopularTitles)
		publicTitle.GET("/recently_updated", h.GetRecentlyUpdatedTitles)
		publicTitle.GET("/new", h.GetNewTitles)
		publicTitle.GET("/:title", h.GetTitle)
	}
}
