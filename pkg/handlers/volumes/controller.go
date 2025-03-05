package volumes

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

	volumesCoversCollection := client.Database("mangacage").Collection("volumes_covers")

	h := handler{
		DB:         db,
		Collection: volumesCoversCollection,
	}

	privateVolume := r.Group("/:title")
	privateVolume.Use(middlewares.AuthMiddleware(secretKey))

	privateVolume.POST("/", h.CreateVolume)
}
