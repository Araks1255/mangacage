package participants

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) ExcludeParticipant(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	participantID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id участника"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesParticipantExist bool

	if err := tx.Raw(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?) AND id != ?) ",
		participantID, claims.ID, claims.ID,
	).Scan(&doesParticipantExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesParticipantExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "участник вашей команды с таким id не найден"})
		return
	}

	result := tx.Exec("UPDATE users SET team_id = NULL WHERE id = ?", participantID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "не удалось исключить участника из вашей команды"})
		return
	}

	result = tx.Exec(
		`DELETE FROM user_roles AS ur
		USING roles AS r WHERE ur.role_id = r.id
		AND ur.user_id = ? AND r.type = 'team'`,
		participantID,
	)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "участник успешно исключен из вашей команды"})
}
