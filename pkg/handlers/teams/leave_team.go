package teams

import (
	"database/sql"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) LeaveTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userTeamID sql.NullInt64
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)

	if !userTeamID.Valid {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы итак не состоите в команде перевода"})
		return
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var userRoles []string
	tx.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if slices.Contains(userRoles, "team_leader") { // Если юзер лидер команды
		if result := tx.Exec( // Назначаем рандомному участнику его команды статус лидера команды
			`UPDATE user_roles SET role_id =
			(SELECT id FROM roles WHERE name = 'team_leader')
			WHERE user_id =
			(
				SELECT id FROM users WHERE team_id =
				(SELECT team_id FROM users WHERE id = ?)
				LIMIT 1
			)`, claims.ID,
		); result.Error != nil {
			log.Println(result.Error)
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
			return
		}
	}

	if result := tx.Exec("UPDATE users SET team_id = null WHERE id = ?", claims.ID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result := tx.Exec(
		`DELETE FROM user_roles AS ur
		INNER JOIN roles AS r ON ur.role_id = r.id
		WHERE ur.user_id = ?
		AND r.type = 'team'`, claims.ID,
	); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "вы успешно покинули команду перевода"})
}
