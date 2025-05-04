package titles

import (
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"

	"github.com/gin-gonic/gin"
)

func (h handler) TranslateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для взятия тайтла на перевод"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingTitleID, titleTeamID uint

	row := tx.Raw("SELECT id, team_id FROM titles WHERE id = ?", desiredTitleID).Row()

	if err = row.Scan(&existingTitleID, &titleTeamID); err != nil {
		log.Println(err)
	}

	if existingTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	if titleTeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл уже переводит другая команда"})
		return
	}

	var userTeamID uint
	tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"}) // По идее это невозможно, ведь сверху идёт проверка на тим лидера, но на практике бд такого не исключает
		return
	}

	if result := tx.Exec("UPDATE titles SET team_id = ? WHERE id = ?", userTeamID, existingTitleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "Теперь ваша команда переводит этот тайтл"})
}
