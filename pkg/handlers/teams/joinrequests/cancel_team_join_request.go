package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CancelTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id заявки должен быть числом"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var existingRequestID uint
	tx.Raw("SELECT id FROM team_join_requests WHERE id = ? AND candidate_id = ?", desiredRequestID, claims.ID).Scan(&existingRequestID)
	if existingRequestID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка не найдена"})
		return
	}

	if result := h.DB.Exec("DELETE FROM team_join_requests WHERE id = ?", existingRequestID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваша заявка на вступление в команду успешно отменена"})
}
