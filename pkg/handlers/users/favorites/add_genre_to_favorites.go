package favorites

import (
	"log"
	"strconv"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) AddGenreToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	genreID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id жанра"})
		return
	}

	err = h.DB.Exec("INSERT INTO user_favorite_genres (user_id, genre_id) VALUES (?, ?)", claims.ID, genreID).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkUserFavoriteGenresGenre) {
			c.AbortWithStatusJSON(404, gin.H{"error": "жанр не найден"})
			return
		}

		if dbErrors.IsUniqueViolation(err, constraints.UserFavoriteGenresPkey) {
			c.AbortWithStatusJSON(409, gin.H{"error": "жанр уже добавлен в ваше избранное"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "жанр успешно добавлен к вам в избранное"})
}
