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

	titlesCoversCollection := client.Database("mangacage").Collection("titles_covers")

	h := handler{
		DB:         db,
		Collection: titlesCoversCollection,
	}

	privateTitle := r.Group("/title")
	privateTitle.Use(middlewares.AuthMiddleware(secretKey))

	privateTitle.POST("/", h.CreateTitle)
	privateTitle.POST("/translate", h.TranslateTitle)
	privateTitle.POST("/:title/subscribe", h.SubscribeToTitle)
}
