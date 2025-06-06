package chapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) DeleteChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesChapterExist bool

	if err := tx.Raw(
		`SELECT EXISTS(
			SELECT 1 FROM chapters AS c
			INNER JOIN volumes AS v ON v.id = c.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
			WHERE c.id = ? AND t.team_id = (SELECT team_id FROM users WHERE id = ?)
		)`, chapterID, claims.ID,
	).Scan(&doesChapterExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesChapterExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена среди глав тайтлов, переводимых вашей командой"})
		return
	}

	result := tx.Exec("DELETE FROM chapters WHERE id = ?", chapterID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить главу"})
		return
	}

	if _, err := h.ChaptersPages.DeleteOne(c.Request.Context(), bson.M{"chapter_id": chapterID}); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "глава успешно удалена"})
}
