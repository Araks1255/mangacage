package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"

	"github.com/gin-gonic/gin"
)

func (h handler) QuitTranslatingTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesTitleExist bool

	if err := tx.Raw(
		"SELECT EXISTS(SELECT 1 FROM titles WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?))",
		titleID, claims.ID,
	).Scan(&doesTitleExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesTitleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди переводимых вашей командой тайтлов"})
		return
	}

	result := tx.Exec("UPDATE titles SET team_id = NULL WHERE id = ?", titleID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при изменении тайтла"})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваша команда больше не переводит этот тайтл"})
}
