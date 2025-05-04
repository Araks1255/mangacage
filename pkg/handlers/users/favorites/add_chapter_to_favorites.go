package favorites

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) AddChapterToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredChapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		UserFavoriteChapterID uint
		ChapterID             uint
	}

	tx.Raw(
		`SELECT
			(SELECT chapter_id FROM user_favorite_chapters WHERE user_id = ? AND chapter_id = ?) AS user_favorite_chapter_id,
			(SELECT id FROM chapters WHERE id = ?) AS chapter_id`,
		claims.ID, desiredChapterID, desiredChapterID,
	).Scan(&existing)

	if existing.UserFavoriteChapterID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "эта глава уже есть у вас в избранном"})
		return
	}
	if existing.ChapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	if result := tx.Exec("INSERT INTO user_favorite_chapters (user_id, chapter_id) VALUES (?, ?)", claims.ID, existing.ChapterID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "глава успешно добавлена к вам в избранное"})
}
