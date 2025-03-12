package chapters

import (
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
	title := strings.ToLower(c.Param("title"))
	volume := strings.ToLower(c.Param("volume"))
	chapterName := strings.ToLower(c.Param("chapter"))

	var chapter models.Chapter
	h.DB.Raw(`SELECT chapters.* FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE chapters.name = ? AND volumes.name = ? AND titles.name = ?`,
		chapterName,
		volume,
		title,
	).Scan(&chapter)

	c.JSON(200, &chapter)
}
