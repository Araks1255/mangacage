package participants

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) AddRoleToParticipant(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userTeamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID == 0 {
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
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для добавления ролей участникам команды"})
		return
	}

	desiredParticipantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id участника команды должен быть числом"})
		return
	}

	var requestBody struct {
		DesiredRoleID uint `json:"roleId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingParticipantID uint
	tx.Raw("SELECT id FROM users WHERE id = ? AND team_id = ?", desiredParticipantID, userTeamID).Scan(&existingParticipantID)
	if existingParticipantID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "участник команды не найден"})
		return
	}

	var existingRoleID uint
	tx.Raw("SELECT id FROM roles WHERE id = ? AND type = 'team'", requestBody.DesiredRoleID).Scan(&existingRoleID)
	if existingRoleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "роль не найдена"})
		return
	}

	var participantRolesIDs []uint
	tx.Raw("SELECT role_id FROM user_roles WHERE user_id = ?", existingParticipantID).Scan(&participantRolesIDs)
	if slices.Contains(participantRolesIDs, existingRoleID) {
		c.AbortWithStatusJSON(409, gin.H{"error": "участник команды уже имеет эту роль"})
		return
	}

	
	if result := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", existingParticipantID, existingRoleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "участнику команды успешно добавлена новая роль"})
}
