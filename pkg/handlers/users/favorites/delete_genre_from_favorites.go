package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteGenreFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var requestBody struct {
		Genre string `json:"genre" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var genreID uint
	h.DB.Raw(`SELECT genres.id FROM genres
		INNER JOIN user_favorite_genres ON user_favorite_genres.genre_id = genres.id
		WHERE user_favorite_genres.user_id = ?
		AND lower(genres.name) = lower(?)`, claims.ID, requestBody.Genre).Scan(&genreID)

	if genreID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "жанр не найден"})
		return
	}

	if result := h.DB.Exec("DELETE FROM user_favorite_genres WHERE user_id = ? AND genre_id = ?", claims.ID, genreID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "жанр успешно удалён из избранного"})
}
