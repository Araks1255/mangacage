package chapters

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
	title := c.Param("title")
	volume := c.Param("volume")
	chapterName := c.Param("chapter")

	var chapter struct {
		ID            uint
		RawCreatedAt  time.Time `json:"-"`
		CreatedAt     string
		Name          string
		Description   string
		NumberOfPages int
		Volume        string
		Title         string
	}

	h.DB.Raw(
		`SELECT c.id, c.created_at AS raw_created_at, c.name, c.description, c.number_of_pages,
		volumes.name AS volume, titles.name AS title
		FROM chapters AS c
		INNER JOIN volumes ON volumes.id = c.volume_id
		INNER JOIN titles ON titles.id = volumes.title_id
		WHERE titles.name = ?
		AND volumes.name = ?
		AND c.name = ?`,
		title, volume, chapterName,
	).Scan(&chapter)

	if chapter.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	chapter.CreatedAt = chapter.RawCreatedAt.Format(time.DateTime)
	c.JSON(200, &chapter)
}
