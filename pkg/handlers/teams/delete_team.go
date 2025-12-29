package teams

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeleteTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	code, err := deleteTeam(h.DB, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "команда успешно удалена"})
}

func deleteTeam(db *gorm.DB, userID uint) (code int, err error) {
	result := db.Exec(
		`WITH user_team_id AS (
			SELECT
				t.id
			FROM
				teams AS t
				INNER JOIN users AS u ON u.team_id = t.id
				INNER JOIN user_roles AS ur ON ur.user_id = u.id
				INNER JOIN roles AS r ON r.id = ur.role_id
			WHERE
				u.id = ? AND r.name = 'team_leader'
		)
		DELETE FROM teams WHERE id = (SELECT id FROM user_team_id)`, // Логика удаления командных ролей вынесена в триггер
		userID,
	)

	if result.Error != nil {
		return 500, err
	}

	if result.RowsAffected == 0 {
		return 404, errors.New("не найдено команды, в которой вы бы были лидером")
	}

	return 0, nil
}
