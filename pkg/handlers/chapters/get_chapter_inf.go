package chapters

import (
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
	title := c.Param("title")
	volume := c.Param("volume")
	chapterName := c.Param("chapter")

	var chapter models.Chapter
	h.DB.Raw(`SELECT chapters.* FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE lower(chapters.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND lower(titles.name) = lower(?)`,
		chapterName,
		volume,
		title,
	).Scan(&chapter)

	c.JSON(200, &chapter)
}
