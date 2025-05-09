package moderation

import (
	"strconv"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyNewChaptersOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10
	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var chapters []models.ChapterOnModerationDTO

	if err := h.DB.Raw(`
		SELECT
			com.id, com.created_at, com.name, com.description, com.number_of_pages,
			v.name AS volume, v.id AS volume_id,
			t.name AS title, t.id AS title_id
		FROM
			chapters_on_moderation AS com
			INNER JOIN volumes AS v ON v.id = com.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			com.existing_id IS NULL
		AND
			com.creator_id = ?
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&chapters).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error":err.Error()})
		return
	}

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших новых глав на модерации"})
		return
	}

	c.JSON(200, &chapters)
}
