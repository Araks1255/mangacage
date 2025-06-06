package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) CancelAppealForProfileChanges(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var profileChangesID *uint

	if err := tx.Raw("DELETE FROM users_on_moderation WHERE existing_id = ? RETURNING id", claims.ID).Scan(&profileChangesID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if profileChangesID == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено изменений вашего профиля на модерации"})
		return
	}

	filter := bson.M{"user_on_moderation_id": *profileChangesID}

	if _, err := h.ProfilePictures.DeleteOne(c.Request.Context(), filter); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "обращение на модерацию изменений профиля успешно отменено"})
}
