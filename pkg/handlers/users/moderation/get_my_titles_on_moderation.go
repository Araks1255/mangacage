package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	titlesHelpers "github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTitlesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10

	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	moderationType := c.Query("type")

	var (
		titles []models.TitleOnModerationDTO
		err    error
	)

	switch moderationType {
	case "new":
		err = titlesHelpers.GetNewTitleOnModeration(h.DB).Where("tom.creator_id = ?", claims.ID).Limit(limit).Scan(&titles).Error

	case "edited":
		err = titlesHelpers.GetEditedTitleOnModeration(h.DB).Where("tom.creator_id = ?", claims.ID).Limit(limit).Scan(&titles).Error

	case "":
		err = titlesHelpers.GetTitleOnModeration(h.DB).Where("tom.creator_id = ?", claims.ID).Limit(limit).Scan(&titles).Error

	default:
		c.AbortWithStatusJSON(400, gin.H{"error": "недопустимый тип модерации"})
		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших тайтлов на модерации"})
		return
	}

	c.JSON(200, &titles)
}
