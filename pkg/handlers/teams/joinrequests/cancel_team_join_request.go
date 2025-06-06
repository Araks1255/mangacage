package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id заявки должен быть числом"})
		return
	}

	result := h.DB.Exec("DELETE FROM team_join_requests WHERE id = ? AND candidate_id = ?", desiredRequestID, claims.ID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "запрос на вступление в команду не найден"}) // Если ошибок нет, но ничего не удалилось, то запись не существует
		return
	}

	c.JSON(200, gin.H{"success": "ваша заявка на вступление в команду успешно отменена"})
}
