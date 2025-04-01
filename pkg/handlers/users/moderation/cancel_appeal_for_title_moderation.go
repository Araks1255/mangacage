package moderation

import (
	"context"
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) CancelAppealForTitleModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var (
		titleID             sql.NullInt64
		titleOnModerationID uint
	)

	row := tx.Raw("SELECT existing_id, id FROM titles_on_moderation WHERE name = ? AND creator_id = ?", title, claims.ID).Row()
	if err := row.Scan(&titleID, &titleOnModerationID); err != nil {
		log.Println(err)
	}

	if titleOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено такого тайтла в ваших заявках на модерацию"})
		return
	}

	if result := tx.Exec("DELETE FROM titles_on_moderation WHERE id = ?", titleOnModerationID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	var filter bson.M

	if titleID.Valid {
		filter = bson.M{"title_id": titleID}
	} else {
		filter = bson.M{"title_on_moderation_id": titleOnModerationID}
	}

	if _, err := h.TitlesCovers.DeleteOne(context.TODO(), filter); err != nil { // По идее, не найденный документ не должен возвращать ошибку
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваша обращение на модерацию отменено"})
}
