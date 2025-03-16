package favorites

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetFavoriteChapters(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var chapters []struct {
		Name string
		Volume string
		Title string
	}

	h.DB.Raw(`SELECT chapters.name, volumes.name AS volume, titles.name AS title
		FROM chapters INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		INNER JOIN user_favorite_chapters ON user_favorite_chapters.chapter_id = chapters.id
		WHERE user_favorite_chapters.user_id = ?
		AND NOT chapters.on_moderation`, claims.ID).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error":"не найдено любимых глав"})
		return
	}

	c.JSON(200, &chapters)
}