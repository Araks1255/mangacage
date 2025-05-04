package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelAppealForVolumeModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredVolumeOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома на модерации"})
		return
	}

	result := h.DB.Exec("DELETE FROM volumes_on_moderation WHERE id = ? AND creator_id = ?", desiredVolumeOnModerationID, claims.ID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден в ваших томах на модерации"})
		return
	}

	c.JSON(200, gin.H{"success": "обращение на модерацию тома успешно отменено"})
}
