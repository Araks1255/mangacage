package chapters

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetVolumeChapters(c *gin.Context) {
	desiredVolumeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тома должен быть числом"})
		return
	}

	var existingVolumeID uint
	h.DB.Raw("SELECT id FROM volumes WHERE id = ?", desiredVolumeID).Scan(&existingVolumeID)
	if existingVolumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	var chapters []models.ChapterDTO

	h.DB.Raw("SELECT id, created_at, name, description, number_of_pages FROM chapters WHERE volume_id = ?", existingVolumeID).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено глав в этом томе"})
		return
	}

	c.JSON(200, &chapters)
}
