package titles

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetTitleCover(c *gin.Context) {
	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var existingTitleID uint
	h.DB.Raw("SELECT id FROM titles WHERE id = ?", desiredTitleID).Scan(&existingTitleID)
	if existingTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	filter := bson.M{"title_id": existingTitleID}

	var result struct {
		TitleID uint   `bson:"title_id"`
		Cover   []byte `bson:"cover"`
	}

	if err := h.TitlesCovers.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Cover) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "обложка тайтла не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
