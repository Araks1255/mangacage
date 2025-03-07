package titles

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetTitleCover(c *gin.Context) {
	title := strings.ToLower(c.Param("title"))

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE name = ?", title).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	filter := bson.M{"title_id": titleID}

	var result struct {
		TitleID uint   `bson:"title_id"`
		Cover   []byte `bson:"cover"`
	}

	if err := h.Collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
