package chapters

import (
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitleChapters(c *gin.Context) {
	title := strings.ToLower(c.Param("title"))

	var chapters []string
	h.DB.Raw("SELECT chapters.name FROM chapters INNER JOIN titles ON chapters.title_id = titles.id WHERE titles.name = ? AND NOT chapters.on_moderation", title).Scan(&chapters)
	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Главы этого тайтла не найдены"})
		return
	}

	response := utils.ConvertToMap(chapters)
	c.JSON(200, &response)
}
