package chapters

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetVolumeChapters(c *gin.Context) {
	volumeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тома должен быть числом"})
		return
	}

	var chapters []struct {
		ID          uint
		Name        string
		Description string
	}

	h.DB.Raw("SELECT id, name, description FROM chapters WHERE volume_id = ?", volumeID).Scan(&chapters)

	c.JSON(200, &chapters)
}
