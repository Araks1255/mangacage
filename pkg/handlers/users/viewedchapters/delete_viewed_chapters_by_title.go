package viewedchapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteViewedChaptersByTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "невалидный id тайтла"})
		return
	}

	result := h.DB.Exec(
		`DELETE FROM
			user_viewed_chapters AS uvc
		USING
			chapters AS c
		WHERE
			c.id = uvc.chapter_id AND
			c.title_id = ? AND uvc.user_id = ?`,
		titleID, claims.ID,
	)

	if result.Error != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено прочитанных вами глав в этом тайтле"})
		return
	}

	c.JSON(200, gin.H{"success": "главы тайтла успешно удалены из вашей истории просмотров"})
}
