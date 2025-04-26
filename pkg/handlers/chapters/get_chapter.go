package chapters

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id главы должен быть числом"})
		return
	}

	var chapter models.ChapterDTO

	h.DB.Raw(
		`SELECT c.id, c.created_at, c.name, c.description, c.number_of_pages,
		v.name AS volume, v.id AS volume_id, t.name AS title, t.id AS title_id
		FROM chapters AS c
		INNER JOIN volumes AS v ON v.id = c.volume_id
		INNER JOIN titles AS t ON t.id = v.title_id
		WHERE c.id = ?`, chapterID,
	).Scan(&chapter)

	if chapter.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	c.JSON(200, &chapter)
}
