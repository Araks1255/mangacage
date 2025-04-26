package titles

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"

	"github.com/gin-gonic/gin"
)

func (h handler) QuitTranslatingTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON user_roles.role_id = roles.id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingTitleID uint
	tx.Raw("SELECT id FROM titles WHERE id = ?", desiredTitleID).Scan(&existingTitleID)
	if existingTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var doesUserTeamTranslatesThisTitle bool
	tx.Raw(
		`SELECT (select team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)`,
		existingTitleID, claims.ID,
	).Scan(&doesUserTeamTranslatesThisTitle)

	if !doesUserTeamTranslatesThisTitle && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	if result := tx.Exec("UPDATE titles SET team_id = null WHERE id = ?", existingTitleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваша команда больше не переводит данный тайтл"})
}
