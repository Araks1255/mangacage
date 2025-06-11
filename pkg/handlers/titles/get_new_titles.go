package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	titlesHelpers "github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	"github.com/gin-gonic/gin"
)

func (h handler) GetNewTitles(c *gin.Context) {
	limit := 10

	if c.Query("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "лимит должен быть числом"})
			return
		}
	}

	var titles []models.TitleDTO

	err := titlesHelpers.GetTitle(h.DB).Order("id DESC").Limit(limit).Scan(&titles).Error
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "новые тайтлы не найдены"})
		return
	}

	c.JSON(200, &titles)
}
