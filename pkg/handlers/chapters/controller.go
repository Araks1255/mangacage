package chapters

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

func RegisterRoutes(db *gorm.DB, collection *mongo.Collection, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	h := handler{
		DB:         db,
		Collection: collection,
	}

	privateChapter := r.Group("/:title/chapter")
	privateChapter.Use(middlewares.AuthMiddleware(secretKey))

	privateChapter.POST("/", h.CreateChapter)
	privateChapter.DELETE(":chapter", h.DeleteChapter)

	r.GET("/get-chapters/:title", h.GetTitleChapters)
	r.GET("/get-chapter/:chapter", h.GetChapter)
	r.GET("/get-chapter-page/:chapter/:page", h.GetChapterPage)
}
