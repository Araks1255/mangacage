package chapters

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetChapterPage(c *gin.Context) {
	title := c.Param("title")
	volume := c.Param("volume")
	chapter := c.Param("chapter")

	numberOfPage, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "номер страницы должен быть числом"})
		return
	}

	var chapterID uint
	h.DB.Raw(`SELECT chapters.id FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE lower(titles.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND lower(chapters.name) = lower(?)`, title, volume, chapter).Scan(&chapterID)
	if chapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	filter := bson.M{"chapter_id": chapterID}

	projection := bson.M{"pages": bson.M{"$slice": []int{numberOfPage, 1}}}

	var result struct {
		Pages [][]byte `bson:"pages"`
	}

	err = h.Collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Pages) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "страница не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.Pages[0])
}
