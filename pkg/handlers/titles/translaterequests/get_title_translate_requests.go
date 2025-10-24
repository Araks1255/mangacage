package translaterequests

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTitleTranslateRequests(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requests []dto.ResponseTitleTranslateRequestDTO

	err := h.DB.Table("titles_translate_requests AS ttr").
		Select("ttr.id, ttr.title_id, t.name AS title").
		Joins("INNER JOIN titles AS t ON t.id = ttr.title_id").
		Where("ttr.team_id = (SELECT team_id FROM users WHERE id = ?)", claims.ID).
		Scan(&requests).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(requests) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено заявок на перевод тайтлов в вашей команде"})
		return
	}

	c.JSON(200, &requests)
}
