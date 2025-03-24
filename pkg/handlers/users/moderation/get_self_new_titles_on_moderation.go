package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h handler) GetSelfNewTitlesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var newTitles []struct {
		CreatedAt   time.Time
		Name        string
		Description string
		Author      string
		Genres      pq.StringArray `gorm:"type:TEXT[]"`
	}

	h.DB.Raw(
		`SELECT t.created_at, t.name, t.description, authors.name AS author, t.genres
		FROM titles_on_moderation AS t
		INNER JOIN authors ON authors.id = t.author_id
		WHERE t.creator_id = ?`, claims.ID,
	).Scan(&newTitles)

	if len(newTitles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших новых тайтлов на модерации"})
		return
	}

	c.JSON(200, &newTitles)
}
