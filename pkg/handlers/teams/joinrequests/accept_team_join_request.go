package joinrequests

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) AcceptTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userTeamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id запроса на вступление в команду должен быть числом"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для приёма заявок на вступление в команду"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var teamJoinRequest models.TeamJoinRequest
	tx.Raw("SELECT * FROM team_join_requests WHERE id = ?", desiredRequestID).Scan(&teamJoinRequest)
	if teamJoinRequest.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "запрос на вступление в команду не найден"})
		return
	}

	if teamJoinRequest.TeamID != userTeamID {
		c.AbortWithStatusJSON(409, gin.H{"error": "запрос на вступление в команду отправлен не в вашу команду"})
		return
	}

	if result := tx.Exec("UPDATE users SET team_id = ? WHERE id = ?", teamJoinRequest.TeamID, teamJoinRequest.CandidateID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	response := make(gin.H, 2)

	if teamJoinRequest.RoleID.Valid {
		if result := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", teamJoinRequest.CandidateID, teamJoinRequest.RoleID); result.Error != nil {
			log.Println(result.Error)
			response["warning"] = "не удалось назначить пользователю роль из запроса"
		}
	}

	tx.Commit()

	response["success"] = "пользователь успешно присоединён к вашей команде"

	c.JSON(200, response)
	// Возможно уведомление юзеру которого приняли сделать
}
