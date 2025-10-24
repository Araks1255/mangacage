package translaterequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTitleTranslateRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	translateRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки на перевод тайтла"})
		return
	}

	var result dto.ResponseTitleTranslateRequestDTO

	err = h.DB.Raw(
		`SELECT
			ttr.*, t.name AS title
		FROM
			titles_translate_requests AS ttr
			INNER JOIN titles AS t ON t.id = ttr.title_id
		WHERE
			ttr.id = ? AND ttr.team_id = (SELECT team_id FROM users WHERE id = ?)`,
		translateRequestID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка на перевод тайтла не найдена среди отправленных вашей командой"})
		return
	}

	c.JSON(200, &result)
}
