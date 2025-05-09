package chapters

import (
	"log"
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

	limit := 10
	if c.Query("limit") != "" {
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var chapters []models.ChapterDTO

	if err := h.DB.Raw(
		"SELECT id, created_at, name, description, number_of_pages FROM chapters WHERE volume_id = ? LIMIT ?",
		desiredVolumeID, limit,
	).Scan(&chapters).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено глав в этом томе"})
		return
	}

	c.JSON(200, &chapters)
}
