package teams

import (
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles on roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if IsUserTeamOwner := slices.Contains(userRoles, "team_owner"); IsUserTeamOwner {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы уже владеете командой перевода"})
		return
	}

	var newTeam models.Team

	if err := c.ShouldBindJSON(&newTeam); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	transaction := h.DB.Begin()

	if result := transaction.Create(&newTeam); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось создать команду"})
		return
	}

	if result := transaction.Exec("UPDATE users SET team_id = ? WHERE id = ?", newTeam.ID, claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось присоеденить вас к команде"})
		return
	}

	if result := transaction.Exec(`INSERT INTO user_roles (user_id, role_id)
		VALUES (?, (SELECT id FROM roles WHERE name = 'team_leader')),
		(?, (SELECT id FROM roles WHERE name = 'translater'))`,
		claims.ID, claims.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось назначить вас лидером команды"})
		return
	}

	transaction.Commit()

	c.JSON(201, gin.H{"success": "Команда успешно создана, и вы являетесь её лидером"})
}
