package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) AddGenreToFavorites(c *gin.Context) {
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
	h.DB.Raw("SELECT id FROM genres WHERE lower(name) = lower(?)", requestBody.Genre).Scan(&genreID)
	if genreID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "жанр не найден"})
		return
	}

	if result := h.DB.Exec("INSERT INTO user_favorite_genres (user_id, genre_id) VALUES (?,?)", claims.ID, genreID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "жанр успешно добавлен в избранное"})
}
