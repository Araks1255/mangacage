package favorites

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, secretKey string, r *gin.Engine) {
	h := handler{DB: db}

	favorites := r.Group("api/users/me/favorites")
	favorites.Use(middlewares.Auth(secretKey))
	{
		favoriteTitles := favorites.Group("/titles")
		{
			favoriteTitles.POST("/:id", h.AddTitleToFavorites)
			favoriteTitles.DELETE("/:id", h.DeleteTitleFromFavorites)
		}

		favoriteChapters := favorites.Group("/chapters")
		{
			favoriteChapters.POST("/:id", h.AddChapterToFavorites)
			favoriteChapters.DELETE("/:id", h.DeleteChapterFromFavorites)
		}

		favoriteGenres := favorites.Group("/genres")
		{
			favoriteGenres.POST("/:id", h.AddGenreToFavorites)
			favoriteGenres.DELETE("/:id", h.DeleteGenreFromFavorites)
		}
	}
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		DB: db,
	}
}
