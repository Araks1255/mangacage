package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyVolumeOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	volumeOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома на модерации"})
		return
	}

	var result dto.ResponseVolumeDTO

	err = h.DB.Raw(
		`SELECT
			vom.*, v.name AS existing, t.name AS title, tom.name AS title_on_moderation
		FROM
			volumes_on_moderation AS vom
			LEFT JOIN volumes AS v ON vom.existing_id = v.id
			LEFT JOIN titles AS t ON vom.title_id = t.id
			LEFT JOIN titles_on_moderation AS tom ON vom.title_on_moderation_id = tom.id
		WHERE
			vom.id = ? AND vom.creator_id = ?`,
		volumeOnModerationID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том с таким id не найден среди ваших заявок на модерацию"})
		return
	}

	c.JSON(200, &result)
}
