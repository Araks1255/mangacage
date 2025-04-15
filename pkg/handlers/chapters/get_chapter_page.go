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
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id главы должен быть числом"})
		return
	}

	numberOfPage, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "номер страницы должен быть числом"})
		return
	}

	filter := bson.M{"chapter_id": chapterID}

	projection := bson.M{"pages": bson.M{"$slice": []int{numberOfPage, 1}}}

	var result struct {
		Pages [][]byte `bson:"pages"`
	}

	err = h.ChaptersPages.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
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
