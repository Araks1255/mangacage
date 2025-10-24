package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTitleOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	var path *string

	err = h.DB.Raw("SELECT cover_path FROM titles_on_moderation WHERE id = ? AND creator_id = ?", titleOnModerationID, claims.ID).Scan(&path).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено обложки тайтла на модерации"})
		return
	}

	c.File(*path)
}
