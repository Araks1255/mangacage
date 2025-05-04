package titles

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

func (h handler) DeleteTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		TitleID             uint
		TitleOnModerationID uint
	}

	tx.Raw(
		`SELECT
			t.id AS title_id, tom.id AS title_on_moderation_id
		FROM 
			titles AS t
			LEFT JOIN titles_on_moderation AS tom ON t.id = tom.existing_id
		WHERE
			t.id = ?`,
		desiredTitleID,
	).Scan(&existing)

	if existing.TitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var titleVolumeID uint
	tx.Raw(`SELECT id FROM volumes WHERE title_id = ? LIMIT 1`, existing.TitleID).Scan(&titleVolumeID)
	if titleVolumeID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "удалить можно только тайтл без томов"})
		return
	}

	var doesUserTeamTranslatesDesiredTitle bool
	tx.Raw(
		"SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)",
		existing.TitleID, claims.ID,
	).Scan(&doesUserTeamTranslatesDesiredTitle)

	if !doesUserTeamTranslatesDesiredTitle && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	if result := tx.Exec("DELETE FROM titles WHERE id = ?", existing.TitleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	filter := bson.M{"title_id": existing.TitleID}

	if _, err := h.TitlesCovers.DeleteOne(context.Background(), filter); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if existing.TitleOnModerationID != 0 {
		filter = bson.M{"title_on_moderation_id": existing.TitleOnModerationID}
		if _, err = h.TitlesOnModerationCovers.DeleteOne(context.Background(), filter); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "тайтл успешно удалён"})

	if _, err := h.TitlesOnModerationCovers.DeleteOne(context.Background(), filter); err != nil {
		log.Println(err)
	}
}
