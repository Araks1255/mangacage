package favorites

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteGenreFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredGenreID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id жанра"})
		return
	}

	result := h.DB.Exec("DELETE FROM user_favorite_genres WHERE user_id = ? AND genre_id = ?", claims.ID, desiredGenreID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "жанр не найден в вашем избранном"})
		return
	}

	c.JSON(200, gin.H{"success": "жанр успешно удален из вашего избранного"})
}
