package participants

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/teams"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) LeaveTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	teamID, teamLeader, newNumberOfTeamParticipants, code, err := leaveTeam(tx, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if teamLeader {
		if newNumberOfTeamParticipants == 0 {
			err = deleteTeam(tx, teamID)
		} else {
			err = teams.TransferTeamLeaderRole(tx, teamID)
		}
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "вы успешно покинули команду перевода"})
}

func leaveTeam(db *gorm.DB, userID uint) (teamID uint, teamLeader bool, newNumberOfTeamParticipants int64, code int, err error) {
	var user struct {
		ID                          *uint
		TeamID                      *uint
		TeamLeader                  bool
		NewNumberOfTeamParticipants *int64
	}

	query :=
		`WITH user_data AS (
			SELECT
				EXISTS(
					SELECT
						1
					FROM
						users AS u
						INNER JOIN user_roles AS ur ON ur.user_id = u.id
						INNER JOIN roles AS r ON r.id = ur.role_id
					WHERE
						u.id = ? AND r.name = 'team_leader'
				) AS team_leader,
				team_id
			FROM
				users
			WHERE
				id = ?
		),
		team AS (
			SELECT
				COUNT(DISTINCT u.id) - 1 AS new_number_of_team_participants 
			FROM
				teams AS t
				INNER JOIN users AS u ON t.id = u.team_id
			WHERE
				u.id = ?
			GROUP BY
				t.id
		)
		UPDATE
			users
		SET
			team_id = NULL
		WHERE
			id = ?
		RETURNING
			id,
			(SELECT team_id FROM user_data), 
			(SELECT team_leader FROM user_data),
			(SELECT new_number_of_team_participants FROM team)` // Логика удаления командных ролей вынесена в триггер

	if err := db.Raw(query, userID, userID, userID, userID).Scan(&user).Error; err != nil {
		return 0, false, 0, 500, err
	}

	if user.ID == nil {
		return 0, false, 0, 409, errors.New("ваш аккаунт не найден в базе данных")
	}

	if user.TeamID == nil {
		return 0, false, 0, 409, errors.New("вы не состоите в команде перевода")
	}

	if user.NewNumberOfTeamParticipants == nil {
		zero := int64(0)
		user.NewNumberOfTeamParticipants = &zero
	}

	return *user.TeamID, user.TeamLeader, *user.NewNumberOfTeamParticipants, 0, nil
}

func deleteTeam(db *gorm.DB, id uint) error {
	return db.Exec("DELETE FROM teams WHERE id = ?", id).Error
}
