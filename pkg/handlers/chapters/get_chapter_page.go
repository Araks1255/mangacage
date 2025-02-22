package chapters

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) GetChapterPages(c *gin.Context) {
	chapter := strings.ToLower(c.Param("chapter"))
	pageNumber, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var pathToChapter string
	h.DB.Raw("SELECT path FROM chapters WHERE name = ?", chapter).Scan(&pathToChapter)
	if pathToChapter == "" {
		c.AbortWithStatusJSON(404, gin.H{"error": "Глава не найдена"})
		return
	}

	pathToPage := fmt.Sprintf("%s/%d.jpg", pathToChapter, pageNumber)

	c.File(pathToPage)
}
