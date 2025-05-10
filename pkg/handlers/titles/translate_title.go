package titles

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"

	"github.com/gin-gonic/gin"
)

func (h handler) TranslateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		DoesTitleExist     bool
		UserTeamID         sql.NullInt64
		IsTitleTranslating bool
	}

	if err := tx.Raw(
		`SELECT
			EXISTS(SELECT 1 FROM titles WHERE id = ?) AS does_title_exist,
			(SELECT team_id FROM users WHERE id = ?) AS user_team_id,
			EXISTS(SELECT 1 FROM titles WHERE id = ? AND team_id IS NOT NULL) AS is_title_translating`,
		titleID, claims.ID, titleID,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !check.DoesTitleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}
	if !check.UserTeamID.Valid {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}
	if check.IsTitleTranslating {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл уже переводится другой командой"})
		return
	}

	result := tx.Exec("UPDATE titles SET team_id = ? WHERE id = ?", check.UserTeamID, titleID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось взять тайтл на перевод"})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "теперь ваша команда переводит этот тайтл"})
}
