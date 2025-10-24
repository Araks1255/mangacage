package participants

import (
	"errors"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeleteParticipantRole(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	ok, isUserTeamLeader, err := checkUserRightsToDeleteParticipantRole(h.DB, claims.ID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для удаления ролей участников команды"})
		return
	}

	participantID, roleID, err := parseDeleteParticipantRoleParams(c.Param, c.ShouldBindJSON)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := checkDeleteParticipantRoleConflicts(tx, claims.ID, participantID, roleID, isUserTeamLeader)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	result := tx.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", participantID, roleID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "участник не обладает такой ролью"})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "участник успешно лишен роли"})
	// Можно участнику уведомление отправить
}

func checkUserRightsToDeleteParticipantRole(db *gorm.DB, userID uint) (ok, isUserTeamLeader bool, err error) {
	var userRoles []string

	err = db.Raw(
		`SELECT
			r.name
		FROM
			roles AS r
			INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE
			ur.user_id = ?`,
		userID,
	).Scan(&userRoles).Error

	if err != nil {
		return false, false, err
	}

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "vice_team_leader") {
		return false, false, nil
	}

	if slices.Contains(userRoles, "team_leader") {
		return true, true, nil
	}

	return true, false, nil
}

func parseDeleteParticipantRoleParams(urlParamFn func(string) string, bindJSONFn func(any) error) (participantID, roleID uint, err error) {
	memberID, err := strconv.ParseUint(urlParamFn("id"), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	var requestBody struct {
		RoleID uint `json:"roleId" binding:"required"`
	}

	if err = bindJSONFn(&requestBody); err != nil {
		return 0, 0, err
	}

	return uint(memberID), requestBody.RoleID, nil
}

func checkDeleteParticipantRoleConflicts(db *gorm.DB, userID, participantID, roleID uint, isUserTeamLeader bool) (code int, err error) {
	var check struct {
		DoesParticipantExist bool
		DoesRoleExist        bool
	}

	var query string

	if isUserTeamLeader {
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team' AND name != 'team_leader') AS does_role_exist`
	} else {
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team' AND name != 'team_leader' AND name != 'vice_team_leader') AS does_role_exist`
	}

	if err = db.Raw(query, participantID, userID, roleID).Scan(&check).Error; err != nil {
		return 500, err
	}

	if !check.DoesParticipantExist {
		return 404, errors.New("участник не найден в вашей команде")
	}

	if !check.DoesRoleExist {
		return 404, errors.New("роль не найдена среди доступных вам для удаления")
	}

	return 0, nil
}
