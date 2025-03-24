package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSelfNewChaptersOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var chapters []struct {
		CreatedAt     time.Time
		Name          string
		Description   string
		NumberOfPages int
		Volume        string
		Title         string
	}

	h.DB.Raw(
		`SELECT c.created_at, c.name, c.description, c.number_of_pages,
		volumes.name AS volume, titles.name AS title
		FROM chapters_on_moderation AS c
		INNER JOIN volumes ON volumes.id = c.volume_id
		INNER JOIN titles ON titles.id = volumes.title_id
		WHERE c.creator_id = ?
		AND c.existing_id IS NULL`, claims.ID,
	).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших новых глав на модерации"})
		return
	}

	c.JSON(200, &chapters)
}
