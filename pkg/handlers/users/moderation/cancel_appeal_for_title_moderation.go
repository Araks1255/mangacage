package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelAppealForTitleModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}

	var titleID uint
	tx.Raw("SELECT id FROM titles_on_moderation WHERE name = ? AND creator_id = ?", requestBody.Name, claims.ID).Scan(&titleID)
	if titleID == 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден в списке ваших тайтлов на модерации"})
		return
	}

	if result := tx.Exec("DELETE FROM titles_on_moderation WHERE id = ?", titleID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваша обращение на модерацию отменено"})
}
