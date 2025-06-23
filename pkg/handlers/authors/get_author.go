package authors

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetAuthor(c *gin.Context) {
	authorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var author models.AuthorDTO

	if err := h.DB.Table("authors").Select("*").Where("id = ?", authorID).Scan(&author).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if author.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "автор не найден"})
		return
	}

	c.JSON(200, &author)
}
