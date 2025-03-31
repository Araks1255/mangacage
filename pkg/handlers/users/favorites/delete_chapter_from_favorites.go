package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteChapterFromFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")
	volume := c.Param("volume")
	chapter := c.Param("chapter")

	var chapterID uint
	h.DB.Raw(
		`SELECT c.id FROM chapters AS c
		INNER JOIN volumes AS v ON c.volume_id = v.id
		INNER JOIN titles AS t ON v.title_id = t.id
		INNER JOIN user_favorite_chapters AS ufc ON ufc.chapter_id = c.id
		WHERE t.name = ? AND v.name = ? AND c.name = ? AND ufc.user_id = ?`,
		title, volume, chapter, claims.ID,
	).Scan(&chapterID)

	if chapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена в ваших избранных главах"})
		return
	}

	if result := h.DB.Exec(
		"DELETE FROM user_favorite_chapters WHERE chapter_id = ? AND user_id = ?",
		chapterID, claims.ID,
	); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "глава успешно удалена из вашего избранного"})
}
