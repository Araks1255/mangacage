package testhelpers

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func CreateTags(db *gorm.DB, quantity int) ([]uint, error) {
	var res []uint

	names := make([]string, quantity)

	for i := 0; i < len(names); i++ {
		names[i] = uuid.New().String()
	}

	err := db.Raw(
		`INSERT INTO tags (name)
		SELECT UNNEST(?::TEXT[])
		RETURNING id`,
		pq.Array(names),
	).Scan(&res).Error

	if err != nil {
		return nil, err
	}

	return res, nil
}
