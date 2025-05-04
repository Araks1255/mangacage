package chapters

import (
	"context"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) DeleteChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для удаления главы"})
		return
	}

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id главы должен быть числом"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var titleID uint
	h.DB.Raw(
		`SELECT t.id FROM titles AS t
		INNER JOIN volumes AS v ON t.id = v.title_id
		INNER JOIN chapters AS c ON v.id = c.volume_id
		WHERE c.id = ?`, chapterID,
	).Scan(&titleID)

	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"}) // Глава не может существовать без тома и тайтла (по ограничению бд, это буквально невозможно), так что если не нашелся тайтл, то и главы такой нет
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	h.DB.Raw(
		"SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)",
		titleID, claims.ID,
	).Scan(&doesUserTeamTranslatesDesiredTitle)

	if !doesUserTeamTranslatesDesiredTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "удалить главу может только команда, выложившая её"})
		return
	}

	if result := tx.Exec("DELETE FROM chapters WHERE id = ?", chapterID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if _, err := h.ChaptersPages.DeleteOne(context.TODO(), bson.M{"chapter_id": chapterID}); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "глава успешно удалена"})
}
