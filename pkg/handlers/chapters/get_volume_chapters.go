package chapters

import (
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) GetVolumeChapters(c *gin.Context) {
	title := c.Param("title")
	volume := c.Param("volume")

	var chapters []string
	h.DB.Raw(
		`SELECT chapters.name FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE volumes.name = lower(?) AND titles.name = lower(?)
		AND NOT chapters.on_moderation`,
		volume,
		title,
	).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "в этом томе ещё нет глав"})
		return
	}

	response := utils.ConvertToMap(chapters)
	c.JSON(200, &response)
}
