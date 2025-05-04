package viewedchapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteViewedChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredChapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "невалидный id главы"})
		return
	}

	result := h.DB.Exec("DELETE FROM user_viewed_chapters WHERE user_id = ? AND chapter_id = ?", claims.ID, desiredChapterID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена среди ваших прочитанных"})
		return
	}

	c.JSON(200, gin.H{"success": "глава успешно удалена из ваших прочитанных"})
}
