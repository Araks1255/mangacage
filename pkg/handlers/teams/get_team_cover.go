package teams

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetTeamCover(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id команды"})
		return
	}

	var path *string

	if err := h.DB.Raw("SELECT cover_path FROM teams WHERE id = ?", teamID).Scan(&path).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	c.File(*path)
}
