package teams

import (
	"errors"

	"gorm.io/gorm"
)

// Использование подразумевается уже после того, как лидер команды больше в ней не состоит
func TransferTeamLeaderRole(db *gorm.DB, teamID uint) error {
	result := db.Exec(
		`INSERT INTO
			user_roles (user_id, role_id)
		SELECT
			id, (SELECT id FROM roles WHERE name = 'team_leader')
		FROM
			users
		WHERE
			team_id = ?
		LIMIT
			1`,
		teamID,
	)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("не удалось назначить другого участника команды на роль лидера")
	}

	return nil
}
