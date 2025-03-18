package volumes

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	h := handler{DB: db}

	privateVolume := r.Group("/volumes/:title")
	privateVolume.Use(middlewares.AuthMiddleware(secretKey))

	privateVolume.POST("/", h.CreateVolume)
	privateVolume.DELETE("/:volume", h.DeleteVolume)

	publicVolume := r.Group("/volumes/:title")
	publicVolume.GET("/", h.GetTitleVolumes)
	publicVolume.GET("/:volume", h.GetVolume)
}
