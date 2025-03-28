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
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON user_roles.role_id = roles.id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if slices.Contains(userRoles, "team_owner") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже являетесь владельцем другой команды"})
		return
	}

	var userTeamID sql.NullInt64
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже состоите в команде перевода"})
		return
	}

	team := c.Param("team")

	var (
		desiredTeamID   uint
		desiredTeamName string
	)

	row := h.DB.Raw("SELECT id, name FROM teams WHERE name = ?", team).Row()

	if err := row.Scan(&desiredTeamID, &desiredTeamName); err != nil {
		log.Println(err)
	}

	if desiredTeamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда перевода не найдена"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if result := tx.Exec("UPDATE users SET team_id = ? WHERE id = ?", desiredTeamID, claims.ID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, (SELECT id FROM roles WHERE name = 'translater'))", claims.ID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": fmt.Sprintf("теперь вы являетесь чатью команды перевода %s", desiredTeamName)})
}
