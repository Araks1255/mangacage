package titles

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	h := handler{DB: db}

	privateTitle := r.Group("/title")
	privateTitle.Use(middlewares.AuthMiddleware(secretKey))

	privateTitle.POST("/", h.CreateTitle)
	privateTitle.POST("/translate", h.TranslateTitle)
	privateTitle.POST("/:title/subscribe", h.SubscribeToTitle)
}
