package favorites

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetFavoriteGenres(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var genres []string
	h.DB.Raw(`SELECT genres.name FROM genres
		INNER JOIN user_favorite_genres ON genres.id = user_favorite_genres.genre_id
		WHERE user_favorite_genres.user_id = ?`, claims.ID).Scan(&genres)

	if len(genres) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено любимых жанров"})
		return
	}

	c.JSON(200, &genres)
}
