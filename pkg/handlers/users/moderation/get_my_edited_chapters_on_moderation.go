package moderation

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyEditedChaptersOnModeration(c *gin.Context) {
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

	h.DB.Raw(
		`SELECT
			com.id, com.created_at, com.name, com.description,
			c.name AS existing, c.id AS existing_id,
			v.name AS volume, v.id AS volume_id,
			t.name AS title, t.id AS title_id
		FROM
			chapters_on_moderation AS com
			INNER JOIN chapters AS c ON c.id = com.existing_id
			INNER JOIN volumes AS v ON v.id = com.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			com.creator_id = ?
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших отредактированных глав на модерации"})
		return
	}

	c.JSON(200, &chapters)
}
