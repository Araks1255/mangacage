package teams

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) GetTeamCover(c *gin.Context) {
	team := c.Param("team")

	var teamID uint
	h.DB.Raw("SELECT id FROM teams WHERE lower(name) = lower(?)", team).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	filter := bson.M{"team_id": teamID}

	var result struct {
		TeamID uint   `bson:"team_id"`
		Cover  []byte `bson:"cover"`
	}

	if err := h.Collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "image/jpeg", result.Cover)
}
