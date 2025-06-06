package moderation

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var allowedEntities = map[string]struct{}{
	"titles":   {},
	"volumes":  {},
	"chapters": {},
	"teams":    {},
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

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var entity struct {
		ID         uint
		ExistingID uint
	}

	query := fmt.Sprintf("DELETE FROM %s_on_moderation WHERE id = ? AND creator_id = ? RETURNING id, existing_id", desiredEntity) // На этом моменте entity уже проверенно, так что её можно подставлять прям так

	if err := tx.Raw(query, entityID, claims.ID).Scan(&entity).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if entity.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено вашей заявки на модерацию"})
		return
	}

	switch desiredEntity {
	case "titles":
		filter := bson.M{"title_on_moderation_id": entityID}
		_, err = h.TitlesCovers.DeleteOne(c.Request.Context(), filter)

	case "teams":
		filter := bson.M{"team_on_moderation_id": entityID}
		_, err = h.TeamsCovers.DeleteOne(c.Request.Context(), filter)

	case "chapters":
		if entity.ExistingID == 0 {
			filter := bson.M{"chapter_on_moderation_id": entityID}

			if res, err := h.ChaptersPages.DeleteOne(c.Request.Context(), filter); res.DeletedCount == 0 {
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при удалении страниц новой главы на модерации"})
				return
			}
		}
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на модерацию успешно отменена"})
}
