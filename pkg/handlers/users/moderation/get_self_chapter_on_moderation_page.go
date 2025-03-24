package moderation

import (
	"context"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetSelfChapterOnModerationPage(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	title := c.Param("title")
	volume := c.Param("volume")
	chapter := c.Param("chapter")

	numberOfPage, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var chapterID uint
	h.DB.Raw(
		`SELECT c.id FROM chapters_on_moderation AS c
		INNER JOIN volumes ON volumes.id = c.volume_id
		INNER JOIN titles ON titles.id = volumes.title_id
		WHERE c.creator_id = ?
		AND titles.name = ?
		AND volumes.name = ?
		AND c.name = ?
		AND c.existing_id IS NULL`,
		claims.ID, title, volume, chapter,
	).Scan(&chapterID)

	if chapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена в ваших главах на модерации"})
		return
	}

	var result struct {
		Pages [][]byte `bson:"pages"`
	}

	filter := bson.M{"chapter_id": chapterID}
	projection := bson.M{"pages": bson.M{"$slice": []int{numberOfPage, 1}}}

	if err = h.ChaptersPages.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Pages) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "страница главы не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.Pages[0])
}
