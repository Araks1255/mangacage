package search

import (
	"github.com/Araks1255/mangacage/pkg/common/models"

	"gorm.io/gorm"
)

func SearchTeams(db *gorm.DB, query string, limit int) (teams *[]models.TeamDTO, err error) {
	var result []models.TeamDTO

	err = db.Raw(
		`SELECT
			t.id, t.created_at, t.name, t.description,
			u.user_name AS leader, u.id AS leader_id
		FROM
			teams AS t
			INNER JOIN users AS u ON t.id = u.team_id
			INNER JOIN user_roles AS ur ON u.id = ur.user_id
			INNER JOIN roles AS r ON ur.role_id = r.id
		WHERE
			r.name = 'team_leader'
		AND
			lower(t.name) ILIKE lower(?)
		LIMIT ?`,
		query, limit,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
