package favorites

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

	privateFavorites := r.Group("api/home/favorites")
	privateFavorites.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateFavoriteTitles := privateFavorites.Group("/titles")
		{
			privateFavoriteTitles.POST("/", h.AddTitleToFavorites)
			privateFavoriteTitles.GET("/", h.GetFavoriteTitles)
			privateFavoriteTitles.DELETE("/:title", h.DeleteTitleFromFavorites)
		}

		privateFavoriteChapters := privateFavorites.Group("/chapters")
		{
			privateFavoriteChapters.POST("/", h.AddChapterToFavorites)
			privateFavoriteChapters.GET("/", h.GetFavoriteChapters)
			privateFavoriteChapters.DELETE("/:title/:volume/:chapter", h.DeleteChapterFromFavorites)
		}

		privateFavoriteGenres := privateFavorites.Group("/genres")
		{
			privateFavoriteGenres.POST("/", h.AddGenreToFavorites)
			privateFavoriteGenres.GET("/", h.GetFavoriteGenres)
			privateFavoriteGenres.DELETE("/:genre", h.DeleteGenreFromFavorites)
		}
	}

}
