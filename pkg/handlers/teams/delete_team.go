package teams

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) DeleteTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		TeamID             *uint
		TeamOnModerationID *uint
	}

	if err := tx.Raw(
		`SELECT
			t.id AS team_id,
			tom.id AS team_on_moderation_id
		FROM
			teams AS t
			LEFT JOIN teams_on_moderation AS tom ON t.id = tom.existing_id
			INNER JOIN users AS u ON u.team_id = t.id
		WHERE
			u.id = ?`,
		claims.ID,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if check.TeamID == nil {
		c.AbortWithStatusJSON(409, gin.H{"error": "ваша команда не найдена"}) // По бизнес логике такого быть не может, так Fкак до этого в middleware идёт проверка ролей пользователя: team_leader и ex_team_leader, и юзер без них до этого момента не дойдет, а если такие роли есть, то команды не быть не может. Но мало ли
		return
	}

	result := h.DB.Exec("DELETE FROM teams WHERE id = ?", *check.TeamID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить команду"})
		return
	}

	filter := bson.M{"team_id": *check.TeamID}

	res, err := h.TeamsCovers.DeleteOne(c.Request.Context(), filter)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось удалить обложку команды"})
		return
	}

	if check.TeamOnModerationID != nil {
		filter = bson.M{"team_on_moderation_id": *check.TeamOnModerationID}

		if _, err = h.TeamsCovers.DeleteOne(c.Request.Context(), filter); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "ваша команда успешно удалена"})
}
