package users

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetViewedChapters(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	type chapter struct {
		ID     uint
		Name   string
		Volume string
		Title  string
	}

	var chapters []chapter
	h.DB.Raw(
		`SELECT chapters.id, chapters.name, volumes.name AS volume, titles.name AS title FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		INNER JOIN user_viewed_chapters ON user_viewed_chapters.chapter_id = chapters.id
		WHERE user_viewed_chapters.user_id = ?
		AND NOT chapters.on_moderation`, claims.ID,
	).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не удалось найти недавно прочитанные вами главы"})
		return
	}

	c.JSON(200, &chapters)
}
