package participants

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
)

func (h handler) AddRoleToParticipant(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	participantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id участника команды должен быть числом"})
		return
	}

	var requestBody struct {
		RoleID uint `json:"roleId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для добавления ролей участникам команды"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		DoesParticipantExist bool
		DoesRoleExist        bool
	}

	var query string
	if slices.Contains(userRoles, "team_leader") {
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team' AND name != 'team_leader') AS does_role_exist`
	} else { // Это чтобы экс тим лидеры не могли назначать других экс тим лидеров (там дополнительная проверка в существовании роли)
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team' AND name != 'team_leader' AND name != 'ex_team_leader') AS does_role_exist`
	}

	if err := tx.Raw(query, participantID, claims.ID, requestBody.RoleID).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !check.DoesParticipantExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "участник вашей команды не найден"})
		return
	}
	if !check.DoesRoleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "роль не найдена среди доступных вам для добавления"})
		return
	}

	err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", participantID, requestBody.RoleID).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UserRolesPkey) {
			c.AbortWithStatusJSON(409, gin.H{"error": "участник команды уже имеет эту роль"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "участнику команды успешно добавлена новая роль"})
	// Уведомление участнику
}
