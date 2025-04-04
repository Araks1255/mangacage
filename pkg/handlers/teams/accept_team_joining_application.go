package teams

import (
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) AcceptTeamJoiningApplication(c *gin.Context) {
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

	candidateName := c.Param("candidate")

	var candidate struct {
		ID              uint
		TeamID          uint
		ApplicationID   uint
		ApplicationRole string
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
		`SELECT u.id, u.team_id, tja.id AS application_id, tja.role AS application_role
		FROM users AS u
		INNER JOIN team_joining_applications AS tja ON tja.candidate_id = u.id
		WHERE u.user_name = ? AND tja.team_id = ?`,
		candidateName, teamID,
	).Scan(&candidate)

	if candidate.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "кандидат не найден"})
		return
	}

	if candidate.ApplicationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "данный пользователь не подавал заявку на вступлнение в вашу команду"})
		return
	}

	if candidate.TeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "кандидат уже является участником другой команды"})
		return
	}

	if result := tx.Exec("UPDATE users SET team_id = ? WHERE user_name = ?", teamID, candidate); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result := tx.Exec(
		"INSERT INTO user_roles (user_id, role_id) VALUES (?, (SELECT id FROM roles WHERE name = ?))",
		claims.ID, candidate.ApplicationRole,
	); result.Error != nil {
		log.Println(result.Error) // Тут ничего страшного, лидер сам поставит если что
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "пользователь успешно присоеденён к вашей команде"})

	if result := h.DB.Exec("DELETE FROM team_joining_applications WHERE candidate_id = ?", claims.ID); result.Error != nil { //  Удаление всех других заявок юзера
		log.Println(result.Error)
	}
	// Возможно уведомление юзеру которого приняли сделать
}
