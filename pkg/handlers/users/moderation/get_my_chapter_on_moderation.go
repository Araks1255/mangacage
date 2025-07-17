package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyChapterOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	var result dto.ResponseChapterDTO

	err = h.DB.Raw(
		`SELECT
			com.*, c.name AS existing,
			t.name AS title, tom.name AS title_on_moderation
		FROM
			chapters_on_moderation AS com
			LEFT JOIN chapters AS c ON com.existing_id = c.id
			LEFT JOIN titles AS t ON com.title_id = t.id
			LEFT JOIN titles_on_moderation AS tom ON com.title_on_moderation_id = tom.id
		WHERE
			com.id = ? AND com.creator_id = ?`,
		chapterOnModerationID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена среди ваших заявок на модерацию"})
		return
	}

	c.JSON(200, &result)
}
