package chapters

import (
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
	desiredChapter := strings.ToLower(c.Param("chapter"))

	var chapter models.Chapter
	h.DB.Raw("SELECT * FROM chapters WHERE name = ?", desiredChapter).Scan(&chapter)
	if chapter.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Глава не найдена"})
		return
	}

	c.JSON(200, &chapter)
}
