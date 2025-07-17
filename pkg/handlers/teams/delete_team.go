package teams

import (
	"context"
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (h handler) DeleteTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	teamID, teamOnModerationID, code, err := getUserTeamsIDs(tx, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := deleteTeam(tx, teamID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := deleteTeamsCovers(c.Request.Context(), h.TeamsCovers, teamID, teamOnModerationID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "команда успешно удалена"})
}

func getUserTeamsIDs(db *gorm.DB, userID uint) (teamID uint, teamOnModerationID *uint, code int, err error) {
	var check struct {
		TeamID             *uint
		TeamOnModerationID *uint
	}

	err = db.Raw(
		`SELECT
			t.id AS team_id,
			tom.id AS team_on_moderation_id
		FROM
			teams AS t
			LEFT JOIN teams_on_moderation AS tom ON t.id = tom.existing_id
			INNER JOIN users AS u ON u.team_id = t.id
		WHERE
			u.id = ?`,
		userID,
	).Scan(&check).Error

	if err != nil {
		return 0, nil, 500, err
	}

	if check.TeamID == nil {
		return 0, nil, 404, errors.New("ваша команда не найдена")
	}

	return *check.TeamID, check.TeamOnModerationID, 0, nil
}

func deleteTeam(db *gorm.DB, teamID uint) error {
	result := db.Exec("DELETE FROM teams WHERE id = ?", teamID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("не удалось удалить вашу команду")
	}

	return nil
}

func deleteTeamsCovers(ctx context.Context, collection *mongo.Collection, teamID uint, teamOnModerationID *uint) error {
	filter := bson.M{"team_id": teamID}

	if _, err := collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	if teamOnModerationID != nil {
		filter = bson.M{"team_on_moderation_id": *teamOnModerationID}

		if _, err := collection.DeleteOne(ctx, filter); err != nil {
			return err
		}
	}

	return nil
}
