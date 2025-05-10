package moderation

import (
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB              *gorm.DB
	TitlesCovers    *mongo.Collection
	ChaptersPages   *mongo.Collection
	ProfilePictures *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, secretKey string, r *gin.Engine) {
	titlesOnModerationCovers := client.Database("mangacage").Collection(mongodb.TitlesOnModerationCoversCollection)
	chaptersOnModerationPages := client.Database("mangacage").Collection(mongodb.ChaptersOnModerationPagesCollection)
	usersOnModerationProfilePictures := client.Database("mangacage").Collection(mongodb.UsersOnModerationProfilePicturesCollection)

	h := handler{
		DB:              db,
		TitlesCovers:    titlesOnModerationCovers,
		ChaptersPages:   chaptersOnModerationPages,
		ProfilePictures: usersOnModerationProfilePictures,
	}

	moderation := r.Group("/api/users/me/moderation")
	moderation.Use(middlewares.Auth(secretKey))
	{
		profile := moderation.Group("/profile")
		{
			profile.GET("/edited", h.GetMyProfileChangesOnModeration)
			profile.GET("/picture", h.GetMyProfilePictureOnModeration)
			profile.DELETE("/edited", h.CancelAppealForProfileChanges)
		}

		titles := moderation.Group("/titles")
		{
			titles.GET("/edited", h.GetMyEditedTitlesOnModeration)
			titles.GET("/new", h.GetMyNewTitlesOnModeration)
			titles.GET("/:id/cover", h.GetMyTitleOnModerationCover)
			titles.DELETE("/:id", h.CancelAppealForTitleModeration)
		}

		chapters := moderation.Group("/chapters")
		{
			chapters.GET("/new", h.GetMyNewChaptersOnModeration)
			chapters.GET("/edited", h.GetMyEditedChaptersOnModeration)
			chapters.GET("/:id/page/:page", h.GetMyChapterOnModerationPage)
			chapters.DELETE("/:id", h.CancelAppealForChapterModeration)
		}

		volumes := moderation.Group("/volumes")
		{
			volumes.GET("/new", h.GetMyNewVolumesOnModeration)
			volumes.GET("/edited", h.GetMyEditedVolumesOnModeration)
			volumes.DELETE("/:id", h.CancelAppealForVolumeModeration)
		}
	}
}

func NewHandler(db *gorm.DB, titlesOnModerationCovers, chaptersOnModerationPages, usersOnModerationProfilePictures *mongo.Collection) handler {
	return handler{
		DB:              db,
		TitlesCovers:    titlesOnModerationCovers,
		ChaptersPages:   chaptersOnModerationPages,
		ProfilePictures: usersOnModerationProfilePictures,
	}
}
