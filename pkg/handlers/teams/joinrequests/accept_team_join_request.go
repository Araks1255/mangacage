package joinrequests

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) AcceptTeamJoinRequest(c *gin.Context) {
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
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для приёма заявок на вступление в команду"})
		return
	}

	desiredRequestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id заявки должен быть числом"})
		return
	}

	var candidate struct {
		ID          uint
		TeamID      uint
		RequestID   uint
		RequestRole string
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	tx.Raw(
		`SELECT u.id, u.team_id, tjr.id AS request_id, tjr.role AS request_role
		FROM team_join_requests AS tjr
		INNER JOIN users AS u ON tjr.candidate_id = u.id
		WHERE tjr.id = ?`, desiredRequestID,
	).Scan(&candidate)

	if candidate.RequestID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "заявка не найдена"})
		return
	}
	if candidate.TeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "кандидат уже является участником другой команды"})
		return
	}

	if result := tx.Exec("UPDATE users SET team_id = ? WHERE id = ?", teamID, candidate.ID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Exec(
		"INSERT INTO user_roles (user_id, role_id) SELECT ?, id FROM roles WHERE name = '' AND type = 'team' ON CONFLICT (user_id, role_id) DO NOTHING",
		candidate.ID, candidate.RequestRole,
	) // Тут ничего страшного, потом можно будет поставить. Да и я подумываю вообще изменить эту систему

	if result := tx.Exec("DELETE FROM team_join_requests WHERE candidate_id = ?", claims.ID); result.Error != nil { //  Удаление всех других заявок юзера
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "пользователь успешно присоеденён к вашей команде"})
	// Возможно уведомление юзеру которого приняли сделать
}
