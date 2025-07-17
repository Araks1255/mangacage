package participants

import (
	"errors"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) AddRoleToParticipant(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	ok, isUserTeamLeader, err := checkUserRightsToAddRoleToParticipant(h.DB, claims.ID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для добавления ролей участникам команды"})
		return
	}

	participantID, roleID, err := parseAddRoleToParticipantParams(c.Param, c.ShouldBindJSON)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := checkAddRoleToParticipantConflicts(tx, claims.ID, participantID, roleID, isUserTeamLeader)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", participantID, roleID).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UserRolesPkey) {
			c.AbortWithStatusJSON(409, gin.H{"error": "участник команды уже имеет эту роль"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	ok, err = didUserTransferTeamLeaderRole(tx, isUserTeamLeader, roleID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if ok {
		result := tx.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = (SELECT id FROM roles WHERE name = 'team_leader')", claims.ID)

		if result.Error != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
			return
		}

		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(500, gin.H{"error": "не удалось снять вас с роли лидера команды при передаче прав"})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "участнику команды успешно добавлена новая роль"})
	// Уведомление участнику
}

func checkUserRightsToAddRoleToParticipant(db *gorm.DB, userID uint) (ok, isUserTeamLeader bool, err error) {
	var userRoles []string

	err = db.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, userID,
	).Scan(&userRoles).Error

	if err != nil {
		return false, false, err
	}

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		return false, false, nil
	}

	if slices.Contains(userRoles, "team_leader") {
		return true, true, nil
	}

	return true, false, nil
}

func parseAddRoleToParticipantParams(urlParamFn func(string) string, bindJSONFn func(any) error) (participantID, roleID uint, err error) {
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

func checkAddRoleToParticipantConflicts(db *gorm.DB, userID, participantID, roleID uint, isUserTeamLeader bool) (code int, err error) {
	var check struct {
		DoesParticipantExist bool
		DoesRoleExist        bool
	}

	var query string

	if isUserTeamLeader {
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team') AS does_role_exist`
	} else { // Это чтобы экс тим лидеры не могли назначать других экс тим лидеров (там дополнительная проверка в существовании роли)
		query = `SELECT
					EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_participant_exist,
					EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team' AND name != 'team_leader' AND name != 'ex_team_leader') AS does_role_exist`
	}

	if err = db.Raw(query, participantID, userID, roleID).Scan(&check).Error; err != nil {
		return 500, err
	}

	if !check.DoesParticipantExist {
		return 404, errors.New("участник не найден в вашей команде")
	}

	if !check.DoesRoleExist {
		return 404, errors.New("роль не найдена среди доступных вам для изменения")
	}

	return 0, nil
}

func didUserTransferTeamLeaderRole(db *gorm.DB, isTeamLeader bool, roleID uint) (bool, error) {
	if !isTeamLeader {
		return false, nil
	}

	var res bool

	result := db.Raw("SELECT (SELECT id FROM roles WHERE name = 'team_leader') = ?", roleID).Scan(&res)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, errors.New("произошла неизвестная ошибка")
	}

	return res, nil
}
