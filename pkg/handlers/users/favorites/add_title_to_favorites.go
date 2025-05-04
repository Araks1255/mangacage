package favorites

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) AddTitleToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		UserFavoriteTitleID uint
		TitleID             uint
	}

	tx.Raw(
		`SELECT
			(SELECT title_id FROM user_favorite_titles WHERE user_id = ? AND title_id = ?) AS user_favorite_title_id,
			(SELECT id FROM titles WHERE id = ?) AS title_id`,
		claims.ID, desiredTitleID, desiredTitleID,
	).Scan(&existing)

	if existing.UserFavoriteTitleID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл уже добавлен в ваше избранное"})
		return
	}
	if existing.TitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	if result := tx.Exec("INSERT INTO user_favorite_titles (user_id, title_id) VALUES (?, ?)", claims.ID, existing.TitleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "тайтл успешно добавлен к вам в избранное"})
}
