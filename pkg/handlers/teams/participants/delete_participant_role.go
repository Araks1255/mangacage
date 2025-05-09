package participants

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteParticipantRole(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	participantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id участника"})
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
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для совершения операций над участниками команды"}) // Я это всё потом в middleware отедльные вынесу
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
					EXISTS(SELECT 1 FROM ROLES WHERE id = ? AND type = 'team' AND name != 'team_leader') AS does_role_exist`
	} else { // Это чтобы экс тим лидеры не могли удалять других экс тим лидеров (там ещё одно условие добавляется на проверке роли)
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM ROLES WHERE id = ? AND type = 'team' AND name != 'team_leader' AND name != 'ex_team_leader') AS does_role_exist`
	}

	if err := tx.Raw(query, participantID, claims.ID, requestBody.RoleID).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !check.DoesParticipantExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "участник не найден в вашей команде"})
		return
	}
	if !check.DoesRoleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "роль не найдена среди доступных вам для удаления"})
		return
	}

	result := tx.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", participantID, requestBody.RoleID)

	if result.Error != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "участник не обладает такой ролью"}) // Ошибка не произошла, юзер существует, роль тоже. Значит, не выполниться запрос может только при отсутствии у юзера роли
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "участник успешно лишен роли"})
	// Можно участнику уведомление отправить
}
