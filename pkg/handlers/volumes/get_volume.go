package volumes

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetVolume(c *gin.Context) {
	desiredVolumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома"})
		return
	}

	var volume models.VolumeDTO

	if err := h.DB.Raw(
		`SELECT
			v.id, v.created_at, v.name, v.description,
			t.name AS title, t.id AS title_id
		FROM
			volumes AS v
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			v.id = ?`,
		desiredVolumeID,
	).Scan(&volume).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if volume.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	c.JSON(200, &volume)
}
