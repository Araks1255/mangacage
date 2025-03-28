package chapters

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                        *gorm.DB
	ChaptersOnModerationPages *mongo.Collection
	ChaptersPages             *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	chaptersOnModerationPagesCollection := client.Database("mangacage").Collection("chapters_on_moderation_pages")
	chapterPagesCollection := client.Database("mangacage").Collection("chapters_pages")

	h := handler{
		DB:                        db,
		ChaptersOnModerationPages: chaptersOnModerationPagesCollection,
		ChaptersPages:             chapterPagesCollection,
	}

	privateChapter := r.Group("/chapters")
	privateChapter.Use(middlewares.AuthMiddleware(secretKey))

	privateChapter.POST("/", h.CreateChapter)
	privateChapter.POST("/:title/:volume/:chapter/edit", h.EditChapter)
	privateChapter.DELETE("/:chapter", h.DeleteChapter)

	publicChapter := r.Group("/chapters/:title/:volume")

	publicChapter.GET("/:chapter/inf", h.GetChapter)
	publicChapter.GET("/", h.GetVolumeChapters)

	r.GET("/chapters/id/:id/page/:page", h.GetChapterPage)
}
