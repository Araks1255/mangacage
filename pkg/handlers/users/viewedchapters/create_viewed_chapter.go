package viewedchapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateViewedChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	err = h.DB.Exec(
		`INSERT INTO
			user_viewed_chapters (user_id, chapter_id, created_at)
		VALUES
			(?, ?, NOW())
		ON CONFLICT
			(user_id, chapter_id) DO UPDATE
		SET
			created_at = EXCLUDED.created_at`,
		claims.ID, chapterID,
	).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, nil)
}
