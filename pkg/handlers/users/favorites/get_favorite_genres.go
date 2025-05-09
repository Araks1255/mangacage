package favorites

import (
	"strconv"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetFavoriteGenres(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10

	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var genres []models.GenreDTO

	if err := h.DB.Raw(
		`SELECT
			g.id, g.name
		FROM
			user_favorite_genres AS uvg
			INNER JOIN genres AS g ON g.id = uvg.genre_id
		WHERE
			uvg.user_id = ?
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&genres).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error":err.Error()})
		return
	}

	if len(genres) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших избранных жанров"})
		return
	}

	c.JSON(200, &genres)
}
