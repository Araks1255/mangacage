package volumes

import (
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitleVolumes(c *gin.Context) {
	title := c.Param("title")

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE name = ?", title).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var volumes []struct {
		ID          uint
		Name        string
		Description string
	}

	h.DB.Raw("SELECT id, name, description FROM volumes WHERE title_id = ?", titleID).Scan(&volumes)

	if len(volumes) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено томов тайтла"})
		return
	}

	c.JSON(200, &volumes)
}
