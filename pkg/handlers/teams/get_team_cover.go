package teams

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetTeamCover(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error":" id команды должен быть числом"})
		return
	}

	filter := bson.M{"team_id": teamID}

	var result struct {
		TeamID uint   `bson:"team_id"`
		Cover  []byte `bson:"cover"`
	}

	if err := h.TeamsCovers.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result.Cover) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error":"обложка команды не найдена"})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
