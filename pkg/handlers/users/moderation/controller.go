package moderation

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB              *gorm.DB
	TitlesCovers    *mongo.Collection
	ChaptersPages   *mongo.Collection
	ProfilePictures *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	titlesOnModerationCovers := client.Database("mangacage").Collection("titles_on_moderation_covers")
	chaptersOnModerationPages := client.Database("mangacage").Collection("chapters_on_moderation_pages")
	usersOnModerationProfilePictures := client.Database("mangacage").Collection("users_on_moderation_profile_pictures")

	h := handler{
		DB:              db,
		TitlesCovers:    titlesOnModerationCovers,
		ChaptersPages:   chaptersOnModerationPages,
		ProfilePictures: usersOnModerationProfilePictures,
	}

	moderation := r.Group("/api/home/moderation")
	moderation.Use(middlewares.AuthMiddleware(secretKey))
	{
		profile := moderation.Group("/profile")
		{
			profile.GET("/edited", h.GetSelfProfileChangesOnModeration)
			profile.GET("/picture", h.GetSelfProfilePictureOnModeration)
			profile.DELETE("/edited", h.CancelAppealForProfileChanges)
		}

		titles := moderation.Group("/titles")
		{
			titles.GET("/edited", h.GetSelfEditedTitlesOnModeration)
			titles.GET("/new", h.GetSelfNewTitlesOnModeration)
			titles.GET("/:title/cover", h.GetSelfTitleOnModerationCover)
			titles.DELETE("/:title", h.CancelAppealForTitleModeration)
		}

		chapters := moderation.Group("/chapters")
		{
			chapters.GET("/new", h.GetSelfNewChaptersOnModeration)
			chapters.GET("/edited", h.GetSelfEditedChaptersOnModeration)
			chapters.GET("/:title/:volume/:chapter/:page", h.GetSelfChapterOnModerationPage)
			chapters.DELETE("/:title/:volume/:chapter", h.CancelAppealForChapterModeration)
		}

		volumes := moderation.Group("/volumes")
		{
			volumes.GET("/new", h.GetSelfNewVolumesOnModeration)
			volumes.GET("/edited", h.GetSelfEditedVolumesOnModeration)
			volumes.DELETE("/:title/:volume", h.CancelAppealForVolumeModeration)
		}
	}
}
