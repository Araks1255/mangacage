package titles

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetTitleCover(c *gin.Context) {
	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var path *string

	if err := h.DB.Raw("SELECT cover_path FROM titles WHERE id = ? AND NOT hidden", titleID).Scan(&path).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	c.File(*path)
}
