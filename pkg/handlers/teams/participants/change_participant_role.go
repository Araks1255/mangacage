package participants

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) ChangeParticipantRole(c *gin.Context) {
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
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для изменения ролей участников команды"})
	}

	var requestBody struct {
		CurrentRole string `json:"currentRole" binding:"required"`
		NewRole     string `json:"newRole" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	desiredParticipantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id участника должно быть числом"})
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

	var existingParticipantID uint
	tx.Raw("SELECT id FROM users WHERE id = ? AND team_id = ?", desiredParticipantID, teamID).Scan(&existingParticipantID)
	if existingParticipantID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "участник команды не найден"})
		return
	}

	var currentRoleID, newRoleID uint

	row := tx.Raw(
		"SELECT (SELECT id FROM roles WHERE name = ?), (SELECT id FROM roles WHERE name = ?)",
		requestBody.CurrentRole, requestBody.NewRole,
	).Row()

	if err := row.Scan(&currentRoleID, &newRoleID); err != nil {
		log.Println(err)
	}

	if currentRoleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "текущая роль участника указана неверно"})
		return
	}
	if newRoleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "новая роль участника указана неверно"})
		return
	}

	if result := tx.Exec("UPDATE user_roles SET role_id = ? WHERE user_id = ? AND role_id = ?", newRoleID, existingParticipantID, currentRoleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "роль участника команды успешно изменена"})
}
