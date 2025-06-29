package viewedchapters

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

	viewedChapters := r.Group("/api/users/me/viewed-chapters")
	viewedChapters.Use(middlewares.Auth(secretKey))
	{
		viewedChapters.DELETE("/:id", h.DeleteViewedChapter)
	}
}
