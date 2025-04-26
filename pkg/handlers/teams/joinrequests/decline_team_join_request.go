package joinrequests

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) DeclineTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var teamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id заявки должен быть числом"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingRequestID uint
	tx.Raw("SELECT id FROM team_join_requests WHERE id = ? AND team_id = ?", desiredRequestID, teamID).Scan(&existingRequestID)
	if existingRequestID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка не найдена"})
		return
	}

	if result := tx.Exec("DELETE FROM team_join_requests WHERE id = ?", existingRequestID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на вступление в вашу команду успешно отменена"})
	// Уведомление кандидату
}
