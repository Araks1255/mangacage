package teams

import (
	"database/sql"
	"fmt"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) JoinTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles ON roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if IsUserTeamOwner := slices.Contains(userRoles, "team_owner"); IsUserTeamOwner {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы уже являетесь владельцем другой команды"})
		return
	}

	var userTeamID sql.NullInt64
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы уже состоите в команде перевода"})
		return
	}

	team := c.Param("team")

	var (
		desiredTeamID   uint
		desiredTeamName string
	)

	row := h.DB.Raw("SELECT id, name FROM teams WHERE lower(name) = lower(?)", team).Row()

	if err := row.Scan(&desiredTeamID, &desiredTeamName); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if desiredTeamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Команда перевода не найдена"})
		return
	}

	transaction := h.DB.Begin()

	if result := transaction.Exec("UPDATE users SET team_id = ? WHERE id = ?", desiredTeamID, claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось присоеденить вас к команде перевода"})
		return
	}

	if result := transaction.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, (SELECT id FROM roles WHERE name = 'translater'))", claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось назначить вас переводчиком"})
		return
	}

	transaction.Commit()

	c.JSON(200, gin.H{"success": fmt.Sprintf("теперь вы являетесь чатью команды перевода %s", desiredTeamName)})
}
