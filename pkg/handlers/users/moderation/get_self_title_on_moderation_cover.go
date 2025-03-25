package moderation

import (
	"context"
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetSelfTitleOnModerationCover(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")

	var titleID, titleOnModerationID sql.NullInt64

	row := h.DB.Raw("SELECT existing_id, id FROM titles_on_moderation WHERE name = ? AND creator_id = ?", title, claims.ID).Row()

	if err := row.Scan(&titleID, &titleOnModerationID); err != nil {
		log.Println(err)
	}

	if !titleID.Valid && !titleOnModerationID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден в ваших тайтлах на модерации"})
		return
	}

	var filter bson.M
	if titleID.Valid {
		filter = bson.M{"title_id": titleID.Int64}
	} else {
		filter = bson.M{"title_on_moderation_id": titleOnModerationID.Int64}
	}

	var titleCover struct {
		TitleID             uint   `bson:"title_id"`
		TitleOnModerationID uint   `bson:"title_on_moderation_id"`
		Cover               []byte `bson:"cover"`
	}

	if err := h.TitlesCovers.FindOne(context.TODO(), filter).Decode(&titleCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", titleCover.Cover)
}
