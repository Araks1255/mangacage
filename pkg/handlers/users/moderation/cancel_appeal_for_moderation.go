package moderation

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

var allowedEntities = map[string]struct{}{
	"authors":  {},
	"genres":   {},
	"tags":     {},
}

func (h handler) CancelAppealForModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredEntity := c.Param("entity")

	if _, ok := allowedEntities[desiredEntity]; !ok {
		c.AbortWithStatusJSON(400, gin.H{"error": "недопустимый тип объекта на модерации"})
		return
	}

	entityID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id"})
		return
	}

	query := fmt.Sprintf("DELETE FROM %s_on_moderation WHERE id = ? AND creator_id = ?", desiredEntity)

	result := h.DB.Exec(query, entityID, claims.ID)

	if result.Error != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка на модерацию не найдена среди оставленных вами"})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию успешно отменена"})
}
