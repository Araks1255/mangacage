package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelTitleTranslateRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id запроса на перевод тайтла"})
		return
	}

	result := h.DB.Exec(
		"DELETE FROM title_translate_requests WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)",
		requestID, claims.ID,
	)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "запрос на перевод тайтла не найден"})
		return
	}

	c.JSON(200, gin.H{"success": "запрос на перевод тайтла успешно отменен"})
}
