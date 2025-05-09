package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyEditedVolumesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10
	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var volumes []models.VolumeOnModerationDTO
	if err := h.DB.Raw(`
		SELECT
			vom.id, vom.created_at, vom.name, vom.description,
			v.name AS existing, v.id AS existing_id,
			t.name AS title, t.id AS title_id
		FROM
			volumes_on_moderation AS vom
			INNER JOIN volumes AS v ON v.id = vom.existing_id
			INNER JOIN titles AS t ON t.id = vom.title_id
			WHERE vom.creator_id = ?
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&volumes).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(volumes) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших отредактированных томов на модерации"})
		return
	}

	c.JSON(200, &volumes)
}
