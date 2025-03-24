package moderation

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB            *gorm.DB
	TitlesCovers  *mongo.Collection
	ChaptersPages *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	titlesOnModerationCovers := client.Database("mangacage").Collection("titles_on_moderation_covers")
	chaptersOnModerationPages := client.Database("mangacage").Collection("chapters_on_moderation_pages")

	h := handler{
		DB:            db,
		TitlesCovers:  titlesOnModerationCovers,
		ChaptersPages: chaptersOnModerationPages,
	}

	moderation := r.Group("/home/moderation")
	moderation.Use(middlewares.AuthMiddleware(secretKey))

	moderation.GET("/titles/edited", h.GetSelfEditedTitlesOnModeration)
	moderation.GET("/titles/new", h.GetSelfNewTitlesOnModeration)
	moderation.DELETE("/titles", h.CancelAppealForTitleModeration)

	moderation.GET("/chapters/new", h.GetSelfNewChaptersOnModeration)
	moderation.GET("/chapters/edited", h.GetSelfEditedChaptersOnModeration)
	moderation.DELETE("/chapters/:title/:volume/:chapter", h.CancelAppealForChapterModeration)

	moderation.GET("/volumes/new", h.GetSelfNewVolumesOnModeration)
	moderation.GET("/volumes/edited", h.GetSelfEditedVolumesOnModeration)
	moderation.DELETE("/volumes/:title/:volume", h.CancelAppealForVolumeModeration)

	moderation.GET("/titles/:title/cover", h.GetSelfTitleOnModerationCover)
	moderation.GET("/chapters/:title/:volume/:chapter/:page", h.GetSelfChapterOnModerationPage)
}
