package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitle(c *gin.Context) {
	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	var title models.TitleDTO

	err = titles.GetTitleWithDependencies(h.DB).Where("t.id = ?", desiredTitleID).Scan(&title).Error
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if title.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	c.JSON(200, &title)
}
