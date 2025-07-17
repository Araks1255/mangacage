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
	titlesCovers := client.Database("mangacage").Collection(mongodb.TitlesCoversCollection)
	chaptersPages := client.Database("mangacage").Collection(mongodb.ChaptersPagesCollection)
	usersProfilePictures := client.Database("mangacage").Collection(mongodb.UsersProfilePicturesCollection)
	teamCovers := client.Database("mangacage").Collection(mongodb.TeamsCoversCollection)

	h := handler{
		DB:              db,
		TitlesCovers:    titlesCovers,
		ChaptersPages:   chaptersPages,
		ProfilePictures: usersProfilePictures,
		TeamsCovers:     teamCovers,
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
			titles.GET("/:id", h.GetMyTitleOnModeration)
		}

		chapters := moderation.Group("/chapters")
		{
			chapters.GET("/", h.GetMyChaptersOnModeration)
			chapters.GET("/:id/page/:page", h.GetMyChapterOnModerationPage)
			chapters.GET("/:id", h.GetMyChapterOnModeration)
		}

		team := moderation.Group("/team")
		{
			team.GET("/", h.GetMyTeamOnModeration)
			team.GET("/cover", h.GetMyTeamOnModerationCover)
		}

		moderation.GET("/authors", h.GetMyAuthorsOnModeration)
		moderation.GET("/genres", h.GetMyGenresOnModeration)
		moderation.GET("/tags", h.GetMyTagsOnModeration)
	}
}

func NewHandler(db *gorm.DB, titlesCovers, chaptersPages, usersPictures, teamsCovers *mongo.Collection) handler {
	return handler{
		DB:              db,
		TitlesCovers:    titlesCovers,
		ChaptersPages:   chaptersPages,
		ProfilePictures: usersPictures,
		TeamsCovers:     teamsCovers,
	}
}
