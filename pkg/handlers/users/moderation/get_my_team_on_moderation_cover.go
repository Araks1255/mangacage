package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTeamOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var path *string

	if err := h.DB.Raw("SELECT cover_path FROM teams_on_moderation WHERE creator_id = ?", claims.ID).Scan(&path).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено обложки вашей команды на модерации"})
		return
	}

	c.File(*path)
}
