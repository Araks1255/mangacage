package chapters

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetChapterPage(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id главы должен быть числом"})
		return
	}

	numberOfPage, err := strconv.ParseUint(c.Param("page"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "номер страницы должен быть числом"})
		return
	}

	var path *string

	if err := h.DB.Raw("SELECT path FROM pages WHERE chapter_id = ? AND number = ? AND NOT hidden", chapterID, numberOfPage).Scan(&path).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "страница не найдена"})
		return
	}

	c.File(*path)
}
