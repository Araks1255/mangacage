package chapters

import (
	"context"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) DeleteChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "удалять главы могут только лидеры команд перевода"})
		return
	}

	title := c.Param("title")
	volume := c.Param("volume")
	chapter := c.Param("chapter")

	var titleID, chapterID uint
	row := h.DB.Raw(
		`SELECT titles.id, chapters.id FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE titles.name = ?
		AND volumes.name = ?
		AND chapters.name = ?`,
		title, volume, chapter,
	).Row()

	if err := row.Scan(&titleID, &chapterID); err != nil {
		log.Println(err)
	}

	if chapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
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

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	if result := tx.Exec("DELETE FROM chapters CASCADE WHERE id = ?", chapterID); result.RowsAffected == 0 {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result, err := h.ChaptersPages.DeleteOne(context.TODO(), bson.M{"chapter_id": chapterID}); result.DeletedCount == 0 {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "глава успешно удалена"})
}
