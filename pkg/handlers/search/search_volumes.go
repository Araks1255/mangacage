package search

import (
	"fmt"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

func (h handler) SearchVolumes(query string, limit int) (volumes *[]models.VolumeDTO, err error) {
	var result []models.VolumeDTO

	err = h.DB.Raw(
		`SELECT
			v.id, v.created_at, v.name, v.description,
			t.name AS title, t.id AS title_id
		FROM
			volumes AS v
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			lower(v.name) ILIKE lower(?)
		LIMIT ?`,
		fmt.Sprintf("%%%s%%", query), limit,
	).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
