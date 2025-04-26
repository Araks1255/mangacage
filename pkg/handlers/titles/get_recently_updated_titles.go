package titles

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetRecentlyUpdatedTitles(c *gin.Context) {
	var limit uint64 = 10

	if c.Query("limit") != "" {
		var err error
		limit, err = strconv.ParseUint(c.Query("limit"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "лимит должен быть числом"})
			return
		}
	}

	var titles []models.TitleDTO

	h.DB.Raw("SELECT * FROM get_recently_updated_titles(?)", limit).Scan(&titles) // Функция описана в ./internal/migrations/sql/create_get_recently_updated_titles.sql

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено тайтлов с недавно вышедшими главами"})
		return
	}

	c.JSON(200, &titles)
}
