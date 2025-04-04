package teams

import (
	"database/sql"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) ChangeParticipantRole(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для изменения ролей участников команды"})
	}

	var requestBody struct {
		CurrentRole string `json:"currentRole" binding:"required"`
		NewRole     string `json:"newRole" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	participant := c.Param("participant")

	var participantID, currentRoleID, newRoleID sql.NullInt64

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	row := tx.Raw(
		`SELECT
		(
			SELECT id from users WHERE team_id = 
			(SELECT team_id FROM users WHERE id = ?)
			AND user_name = ?
		),
		(
			SELECT r.id FROM roles AS r
			INNER JOIN user_roles AS ur ON ur.role_id = r.id
			INNER JOIN users AS u ON u.id = ur.user_id
			WHERE u.user_name = ? AND r.name = ?
		),
		(SELECT id FROM roles WHERE name = ?)`,
		claims.ID, participant, participant, requestBody.CurrentRole, requestBody.NewRole,
	).Row()

	if err := row.Scan(&participantID, &currentRoleID, &newRoleID); err != nil {
		log.Println(err)
	}

	if !participantID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "участник команды не найден"})
		return
	}

	if !currentRoleID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "текущая роль участника указана неверно"})
		return
	}
	if !newRoleID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "новая роль указана неверно"})
		return
	}

	if result := tx.Exec("UPDATE user_roles SET role_id = ? WHERE user_id = ? AND role_id = ?", newRoleID, participantID, currentRoleID); result.Error != nil { // !!!!!!!!
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "роль участника команды успешно изменена"})
}
