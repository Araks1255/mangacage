package users

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/teams"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) DeleteProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	userTeamID, userTeamLeader, newNumberOfTeamParticipants, code, err := deleteUser(tx, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if userTeamLeader && userTeamID != nil {
		if newNumberOfTeamParticipants == 0 {
			err = deleteTeam(tx, *userTeamID)
		} else {
			err = teams.TransferTeamLeaderRole(tx, *userTeamID)
		}
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.SetCookie("mangacage_token", "", -1, "/", "localhost", false, true) // ПОМЕНЯТЬ НА ПРОДЕ

	c.JSON(200, gin.H{"success": "ваш аккаунт успешно удален"})
}

func deleteUser(db *gorm.DB, id uint) (teamID *uint, teamLeader bool, newNumberOfTeamParticipants int64, code int, err error) {
	var deletedUser struct {
		ID                          *uint
		TeamID                      *uint
		TeamLeader                  bool
		NewNumberOfTeamParticipants *int64
	}

	query :=
		`WITH deleting_user AS (
			SELECT EXISTS(
				SELECT 
					1
				FROM
					users AS u
					INNER JOIN teams AS t ON t.id = u.team_id
					INNER JOIN user_roles AS ur ON ur.user_id = u.id
					INNER JOIN roles AS r ON r.id = ur.role_id
				WHERE
					u.id = ? AND r.name = 'team_leader'
			) AS team_leader
		)
		DELETE FROM
			users AS u
		WHERE
			u.id = ?
		RETURNING
			u.id, u.team_id,
			(SELECT team_leader FROM deleting_user) AS team_leader,
			(SELECT number_of_participants FROM teams WHERE id = u.team_id) - 1 AS new_number_of_team_participants`

	err = db.Raw(query, id, id).Scan(&deletedUser).Error

	if err != nil {
		return nil, false, 0, 500, err
	}

	if deletedUser.ID == nil {
		return nil, false, 0, 404, errors.New("ваш аккаунт не был найден в базе данных")
	}

	if deletedUser.NewNumberOfTeamParticipants == nil { // Такого быть не должно, ведь в базе это NOT NULL столбец, но если в структуре сделать тип просто int64, то при скане будет ещё одна ошибка, перекрывающая основную (если запрос не выполнится)
		zero := int64(0)
		deletedUser.NewNumberOfTeamParticipants = &zero
	}

	return deletedUser.TeamID, deletedUser.TeamLeader, *deletedUser.NewNumberOfTeamParticipants, 0, nil
}

func deleteTeam(db *gorm.DB, id uint) error {
	return db.Exec("DELETE FROM teams WHERE id = ?", id).Error
}
