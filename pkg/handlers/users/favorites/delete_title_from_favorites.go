package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteTitleFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var requestBody struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var titleID uint
	h.DB.Raw(`SELECT titles.id FROM titles
		INNER JOIN user_favorite_titles ON user_favorite_titles.title_id = titles.id
		WHERE user_favorite_titles.user_id = ?
		AND lower(titles.name) = lower(?)`, claims.ID, requestBody.Title).Scan(&titleID)

	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	if result := h.DB.Exec("DELETE FROM user_favorite_titles WHERE user_id = ? AND title_id = ?", claims.ID, titleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "тайтл успешно удалён из избранного"})
}
