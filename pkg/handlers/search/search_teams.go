package search

import (
	"fmt"
	"time"
)

type Team struct {
	ID          uint
	CreatedAt   time.Time
	Name        string
	Description string
	Leader      string
}

func (h handler) SearchTeams(query string, limit int) (teams *[]Team, quantity int) {
	var result []Team

	h.DB.Raw(
		`SELECT t.id, t.created_at, t.name, t.description, users.user_name AS leader
		FROM teams AS t
		INNER JOIN users ON t.id = users.team_id
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'team_leader'
		AND lower(t.name) ILIKE lower(?)
		LIMIT ?`, fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result)

	return &result, len(result)
}
