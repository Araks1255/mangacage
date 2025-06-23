package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyGenresOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := uint64(10)
	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.ParseUint(c.Query("limit"), 10, 32); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var result []models.GenreOnModerationDTO

	if err := h.DB.Table("genres_on_moderation").Select("*").Where("creator_id = ?", claims.ID).Limit(int(limit)).Scan(&result).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших жанров на модерации"})
		return
	}

	c.JSON(200, &result)
}
