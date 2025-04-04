package teams

import (
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) DeclineTeamJoiningApplication(c *gin.Context) {
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
		INNER JOIN user_roles AS ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	candidate := c.Param("candidate")

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var applicationID uint
	tx.Raw(
		`SELECT tja.id FROM team_joining_applications AS tja
		INNER JOIN users AS u ON u.id = tja.candidate_id
		WHERE tja.team_id = ? AND u.user_name = ?`,
		teamID, candidate,
	).Scan(&applicationID)

	if applicationID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "пользователь с таким именем не подавал заявку на вступление в вашу команду"})
		return
	}

	if result := tx.Exec("DELETE FROM team_joining_applications WHERE id = ?", applicationID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на вступление в вашу команду успешно отменена"})
	// Уведомление кандидату
}
