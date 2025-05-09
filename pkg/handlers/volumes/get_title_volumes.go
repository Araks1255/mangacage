package volumes

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitleVolumes(c *gin.Context) {
	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	limit := 10
	if c.Query("limit") != "" {
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var volumes []models.VolumeDTO

	if err := h.DB.Raw(
		`SELECT
			v.id, v.created_at, v.name, v.description,
			t.name AS title, t.id AS title_id
		FROM
			volumes AS v
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			t.id = ?
		LIMIT ?`,
		titleID, limit,
	).Scan(&volumes).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(volumes) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено томов этого тайтла"})
		return
	}

	c.JSON(200, &volumes)
}
