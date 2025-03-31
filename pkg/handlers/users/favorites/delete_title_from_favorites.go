package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteTitleFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")

	var titleID uint
	h.DB.Raw(
		`SELECT t.id FROM titles AS t
		INNER JOIN user_favorite_titles AS uft ON uft.title_id = t.id
		WHERE uft.user_id = ? AND t.name = ?`, claims.ID, title,
	).Scan(&titleID)

	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	if result := h.DB.Exec(
		"DELETE FROM user_favorite_titles WHERE user_id = ? AND title_id = ?",
		claims.ID, titleID,
	); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "тайтл успешно удалён из избранного"})
}
