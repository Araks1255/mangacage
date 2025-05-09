package viewedchapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetUserViewedChapters(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10

	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var chapters []models.ChapterDTO

	if err := h.DB.Raw(
		`SELECT
			c.id, c.name, uvc.created_at,
			v.name AS volume, v.id AS volume_id,
			t.name AS title, t.id AS title_id
		FROM
			user_viewed_chapters AS uvc
			INNER JOIN chapters AS c ON c.id = uvc.chapter_id
			INNER JOIN volumes  AS v ON v.id = c.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
			WHERE uvs.user_id = ?
		ORDER BY
			uvc.created_at DESC
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&chapters).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас ещё нет прочитанных глав"})
		return
	}

	c.JSON(200, &chapters)
}
