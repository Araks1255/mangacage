package titles

import (
	"database/sql"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) DeleteTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		DoesTitleExist      bool
		TitleOnModerationID sql.NullInt64
	}

	err = tx.Raw(
		`SELECT
			EXISTS(SELECT 1 FROM titles WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS does_title_exist,
			(SELECT id FROM titles_on_moderation WHERE existing_id = ?) AS title_on_moderation_id`,
		titleID, claims.ID, titleID,
	).Scan(&check).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !check.DoesTitleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди переводимых вашей командой тайтлов"})
		return
	}

	result := tx.Exec("DELETE FROM titles WHERE id = ?", titleID)

	if result.Error != nil {
		if dbErrors.IsForeignKeyViolation(result.Error, constraints.FkVolumesTitle) {
			c.AbortWithStatusJSON(409, gin.H{"error": "удалить можно только тайтл без томов"})
		} else {
			log.Println(result.Error)
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при удалении тайтла"})
		return
	}

	filter := bson.M{"title_id": titleID}

	mongoResult, err := h.TitlesCovers.DeleteOne(c.Request.Context(), filter)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if mongoResult.DeletedCount == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при удалении обложки тайтла"})
		return
	}

	if check.TitleOnModerationID.Valid {
		filter := bson.M{"title_on_moderation_id": check.TitleOnModerationID.Int64}
		if _, err := h.TitlesOnModerationCovers.DeleteOne(c.Request.Context(), filter); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "тайтл успешно удалён"})
}
