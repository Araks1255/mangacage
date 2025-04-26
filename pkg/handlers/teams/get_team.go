package teams

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id команды должен быть числом"})
		return
	}

	var team struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	h.DB.Raw("SELECT id, name, description FROM teams WHERE id = ?", teamID).Scan(&team)
	if team.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	c.JSON(200, &team)
}
