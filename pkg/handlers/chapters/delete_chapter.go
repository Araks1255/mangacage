package chapters

import (
	"context"
	"log"
	"slices"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) DeleteChapter(c *gin.Context) { // НАЗВАНИЕ ГЛАВЫ НЕУНИКАЛЬНО. ПЕРЕДЕЛАТЬ
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles ON roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if isUserTeamLeader := slices.Contains(userRoles, "team_leader"); !isUserTeamLeader {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы не являетесь лидером команды перевода"})
		return
	}

	title := strings.ToLower(c.Param("title"))
	volume := strings.ToLower(c.Param("volume"))
	chapter := strings.ToLower(c.Param("chapter"))

	var titleID, chapterID uint
	row := h.DB.Raw(
		"SELECT titles.id, chapters.id FROM chapters INNER JOIN volumes ON chapters.volume_id = volumes.id INNER JOIN titles ON volumes.title_id = titles.id WHERE titles.name = ? AND volumes.name = ? AND chapters.name = ?",
		title,
		volume,
		chapter).Row()

	if err := row.Scan(&titleID, &chapterID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var WasDesiredChapterCreatedByUserTeam bool
	h.DB.Raw(
		"SELECT CAST(CASE WHEN (SELECT team_id FROM users WHERE id = ?) = (SELECT team_id FROM titles WHERE id = ?) THEN TRUE ELSE FALSE END AS BOOLEAN)",
		claims.ID,
		titleID).Scan(&WasDesiredChapterCreatedByUserTeam)

	if !WasDesiredChapterCreatedByUserTeam {
		c.AbortWithStatusJSON(403, gin.H{"error": "Удалить главу может только команда, выложившая её"})
		return
	}

	tx := h.DB.Begin()

	if result := tx.Exec("DELETE FROM chapters WHERE id = ?", chapterID); result.RowsAffected == 0 {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result, err := h.Collection.DeleteOne(context.TODO(), bson.M{"chapter_id": chapterID}); result.DeletedCount == 0 {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "Глава успешно удалена"})
}
