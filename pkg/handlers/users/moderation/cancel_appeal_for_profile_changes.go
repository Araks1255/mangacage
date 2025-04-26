package moderation

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) CancelAppealForProfileChanges(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userOnModerationID uint
	h.DB.Raw("SELECT id  FROM users_on_moderation WHERE existing_id = ?", claims.ID).Scan(&userOnModerationID)
	if userOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "у вас нет изменений профиля, ожидающих модерации"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := h.DB.Exec("DELETE FROM users_on_moderation WHERE id = ?", userOnModerationID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	filter := bson.M{"user_id": claims.ID}
	if _, err := h.ProfilePictures.DeleteOne(context.TODO(), filter); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "обращение на модерацию изменений профиля успешно отменено"})
}
