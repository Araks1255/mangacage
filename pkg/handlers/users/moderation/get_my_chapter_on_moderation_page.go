package moderation

import (
	"context"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetMyChapterOnModerationPage(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredChapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	numberOfPage, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный номер страницы"})
		return
	}

	var existingChapterOnModerationID uint
	h.DB.Raw("SELECT id FROM chapters_on_moderation WHERE id = ? AND creator_id = ?", desiredChapterOnModerationID, claims.ID).Scan(&existingChapterOnModerationID)
	if existingChapterOnModerationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена в списке ваших глав на модерации"})
		return
	}

	filter := bson.M{"chapter_on_moderation_id": existingChapterOnModerationID}
	projection := bson.M{"pages": bson.M{"$slice": []int{numberOfPage, 1}}}

	var result struct {
		Pages [][]byte `bson:"pages"`
	}

	if err = h.ChaptersPages.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Pages) == 0 || result.Pages[0] == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "страница главы на модерации не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.Pages[0])
}
