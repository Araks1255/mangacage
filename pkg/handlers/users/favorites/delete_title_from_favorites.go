package favorites

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteTitleFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	result := h.DB.Exec("DELETE FROM user_favorite_titles WHERE user_id = ? AND title_id = ?", claims.ID, desiredTitleID)

	if result.Error != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден в вашем избранном"})
		return
	}

	c.JSON(200, gin.H{"success": "тайтл успешно удалён из вашего избранного"})
}
