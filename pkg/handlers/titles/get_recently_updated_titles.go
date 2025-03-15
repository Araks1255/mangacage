package titles

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) GetRecentlyUpdatedTitles(c *gin.Context) {
	limit, err := strconv.Atoi(c.Param("limit"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error":err.Error()})
		return
	}

	var titles []string
	h.DB.Raw(`SELECT titles.name FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE NOT titles.on_moderation
		ORDER BY chapters.updated_at DESC
		LIMIT ?`, limit).Scan(&titles)

	if len(titles) == 0 { // Ну мало ли
		c.AbortWithStatusJSON(404, gin.H{"error":"не найдено недавно обновлённых тайтлов"})
		return
	}

	response := utils.ConvertToMap(titles)
	c.JSON(200, &response)
}