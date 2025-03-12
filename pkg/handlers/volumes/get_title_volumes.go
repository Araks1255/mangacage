package volumes

import (
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitleVolumes(c *gin.Context) {
	title := strings.ToLower(c.Param("title"))

	var volumes []string
	h.DB.Raw(
		`SELECT volumes.name FROM volumes
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE titles.name = ?
		AND NOT volumes.on_moderation`,
		title,
	).Scan(&volumes)

	if len(volumes) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "в этом тайтле ещё нет томов"})
		return
	}

	response := utils.ConvertToMap(volumes)
	c.JSON(200, &response)
}
