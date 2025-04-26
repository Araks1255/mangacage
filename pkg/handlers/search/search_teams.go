package search

import (
	"fmt"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

func (h handler) SearchTeams(query string, limit int) (teams *[]models.TeamDTO, quantity int) {
	var result []models.TeamDTO

	h.DB.Raw(
		`SELECT t.id, t.created_at, t.name, t.description,
		u.user_name AS leader, u.id AS leader_id
		FROM teams AS t
		INNER JOIN users AS u ON t.id = u.team_id
		INNER JOIN user_roles AS ur ON u.id = ur.user_id
		INNER JOIN roles AS r ON ur.role_id = r.id
		WHERE r.name = 'team_leader'
		AND lower(t.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
