package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteGenreFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	genre := c.Param("genre")

	var genreID uint
	h.DB.Raw(
		`SELECT g.id FROM genres AS g
		INNER JOIN user_favorite_genres AS ufg ON ufg.genre_id = g.id
		WHERE ufg.user_id = ? AND g.name = ?`,
		claims.ID, genre,
	).Scan(&genreID)

	if genreID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "жанр не найден"})
		return
	}

	if result := h.DB.Exec(
		"DELETE FROM user_favorite_genres WHERE user_id = ? AND genre_id = ?",
		claims.ID, genreID,
	); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "жанр успешно удалён из избранного"})
}
