package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSelfNewVolumesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var volumes []struct {
		CreatedAt   time.Time
		Name        string
		Description string
		Title       string
	}

	h.DB.Raw(
		`SELECT v.created_at, v.name, v.description, titles.name AS title
		FROM volumes_on_moderation AS v
		INNER JOIN titles ON titles.id = v.title_id
		WHERE v.creator_id = ?
		AND existing_id IS NULL`,
		claims.ID,
	).Scan(&volumes)

	if len(volumes) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших новых томов на модерации"})
		return
	}

	c.JSON(200, &volumes)
}
