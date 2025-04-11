package chapters

import (
	"github.com/gin-gonic/gin"
)

func (h handler) GetVolumeChapters(c *gin.Context) {
	title := c.Param("title")
	volume := c.Param("volume")

	var chapters []struct {
		Name        string
		Description string
	}

	h.DB.Raw(
		`SELECT c.name, c.description FROM chapters AS c
		INNER JOIN volumes AS v on v.id = c.volume_id
		INNER JOIN titles AS t ON t.id = v.title_id
		WHERE t.name = ? AND v.name = ?`,
		title, volume,
	).Scan(&chapters)

	c.JSON(200, &chapters)
}
