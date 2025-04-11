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

	privateChapter := r.Group("/api/chapters/:title/:volume")
	privateChapter.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateChapter.POST("/", h.CreateChapter)
		privateChapter.DELETE("/:chapter", h.DeleteChapter)
		privateChapter.POST("/:chapter/edited", h.EditChapter) // Тут идёт создание отредактированной главы (прямо отдельная сущность в отдельной таблице базы данных), поэтому post а не put
	}

	publicChapter := r.Group("/api/chapters/:title/:volume")
	{
		publicChapter.GET("/:chapter", h.GetChapter)
		publicChapter.GET("/", h.GetVolumeChapters)
	}

	r.GET("/api/chapters/id/:id/page/:page", h.GetChapterPage)
}
