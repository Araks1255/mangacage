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
	TeamsCovers     *mongo.Collection
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, secretKey string, r *gin.Engine) {
	titlesOnModerationCovers := client.Database("mangacage").Collection(mongodb.TitlesOnModerationCoversCollection)
	chaptersOnModerationPages := client.Database("mangacage").Collection(mongodb.ChaptersOnModerationPagesCollection)
	usersOnModerationProfilePictures := client.Database("mangacage").Collection(mongodb.UsersOnModerationProfilePicturesCollection)
	teamsOnModerationCovers := client.Database("mangacage").Collection(mongodb.TeamsOnModerationCoversCollection)

	h := handler{
		DB:              db,
		TitlesCovers:    titlesOnModerationCovers,
		ChaptersPages:   chaptersOnModerationPages,
		ProfilePictures: usersOnModerationProfilePictures,
		TeamsCovers:     teamsOnModerationCovers,
	}

	moderation := r.Group("/api/users/me/moderation")
	moderation.Use(middlewares.Auth(secretKey))
	{
		moderation.DELETE("/:entity/:id", h.CancelAppealForModeration)

		profile := moderation.Group("/profile")
		{
			profile.GET("/edited", h.GetMyProfileChangesOnModeration)
			profile.GET("/picture", h.GetMyProfilePictureOnModeration)
			profile.DELETE("/edited", h.CancelAppealForProfileChanges)
		}

		titles := moderation.Group("/titles")
		{
			titles.GET("/", h.GetMyTitlesOnModeration)
			titles.GET("/:id/cover", h.GetMyTitleOnModerationCover)
		}

		chapters := moderation.Group("/chapters")
		{
			chapters.GET("/", h.GetMyChaptersOnModeration)
			chapters.GET("/:id/page/:page", h.GetMyChapterOnModerationPage)
		}

		volumes := moderation.Group("/volumes")
		{
			volumes.GET("/", h.GetMyVolumesOnModeration)
		}

		teams := moderation.Group("/teams")
		{
			teams.GET("/team", h.GetMyTeamOnModeration)
		}
	}
}

func NewHandler(db *gorm.DB, titlesOnModerationCovers, chaptersOnModerationPages, usersOnModerationProfilePictures, teamsOnModerationCovers *mongo.Collection) handler {
	return handler{
		DB:              db,
		TitlesCovers:    titlesOnModerationCovers,
		ChaptersPages:   chaptersOnModerationPages,
		ProfilePictures: usersOnModerationProfilePictures,
		TeamsCovers:     teamsOnModerationCovers,
	}
}
