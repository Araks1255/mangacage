package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) UnSubscribeFromTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	result := h.DB.Exec("DELETE FROM user_titles_subscribed_to WHERE user_id = ? AND title_id = ?", claims.ID, titleID)

	if result.Error != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди ваших подписок"})
		return
	}

	c.JSON(200, gin.H{"success": "вы успешно отписались от тайтла"})
}
