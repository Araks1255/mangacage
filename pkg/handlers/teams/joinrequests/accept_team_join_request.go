package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
)

func (h handler) AcceptTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id заявки на вступление в команду"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var teamJoinRequest models.TeamJoinRequest

	if err := tx.Raw(
		"SELECT * FROM team_join_requests WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)",
		desiredRequestID, claims.ID,
	).Scan(&teamJoinRequest).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if teamJoinRequest.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка не найдена среди заявок на вступление в вашу команду"})
		return
	}

	result := tx.Exec("UPDATE users SET team_id = ? WHERE id = ?", teamJoinRequest.TeamID, teamJoinRequest.CandidateID)

	if result.Error != nil {
		log.Println(result.Error.Error())
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось присоединить кандидата к команде"})
		return
	}

	response := make(gin.H, 2)

	if teamJoinRequest.RoleID.Valid {
		err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", teamJoinRequest.CandidateID, teamJoinRequest.RoleID.Int64).Error

		if err != nil {
			if dbErrors.IsUniqueViolation(err, constraints.UserRolesPkey) {
				response["warning"] = "не удалось назначить роль из заявки кандидату (кандидат уже имеет такую роль)"
			} else if dbErrors.IsForeignKeyViolation(err, constraints.FkUserRolesRole) {
				response["warning"] = "не удалось назначить роль из заявки кандидату (указан id несуществующей роли)"
			} else {
				log.Println(err)
				response["warning"] = "не удалось назначить роль из заявки кандидату по неизвестной причине"
			}
		}
	}

	if err := tx.Exec("DELETE FROM team_join_requests WHERE candidate_id = ?", teamJoinRequest.CandidateID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	response["success"] = "пользователь успешно присоединён к вашей команде"

	c.JSON(200, response)
	// Уведомление кандидату
}
