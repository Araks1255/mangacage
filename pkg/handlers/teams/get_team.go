package teams

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id команды"})
		return
	}

	var team models.TeamDTO

	if err = h.DB.Table("teams").Select("*").Where("id = ?", teamID).Scan(&team).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if team.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	c.JSON(200, &team)
}
