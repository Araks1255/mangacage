package favorites

import (
	"log"
	"strconv"

	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) AddGenreToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredGenreID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id жанра"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		UserFavoriteGenreID uint
		GenreID             uint
	}

	tx.Raw(
		`SELECT
			(SELECT genre_id FROM user_favorite_genres WHERE user_id = ? AND genre_id = ?) AS user_favorite_genre_id,
			(SELECT id FROM genres WHERE id = ?) AS genre_id`,
		claims.ID, desiredGenreID, desiredGenreID,
	).Scan(&existing)

	if existing.UserFavoriteGenreID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "этот жанр уже добавлен к вам в избранное"})
		return
	}
	if existing.GenreID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "жанр не найден"})
		return
	}

	if result := tx.Exec("INSERT INTO user_favorite_genres (user_id, genre_id) VALUES (?, ?)", claims.ID, existing.GenreID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "жанр успешно добавлен к вам в избранное"})
}
