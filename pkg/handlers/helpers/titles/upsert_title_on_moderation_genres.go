package titles

import (
	"errors"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func UpsertTitleOnModerationGenres(db *gorm.DB, id uint, genresIDs []uint) (code int, err error) {
	if len(genresIDs) == 0 {
		return 0, nil
	}

	if err := db.Exec("DELETE FROM title_on_moderation_genres WHERE title_on_moderation_id = ?", id).Error; err != nil {
		return 500, err
	}

	err = db.Exec(
		"INSERT INTO title_on_moderation_genres (title_on_moderation_id, genre_id) SELECT ?, UNNEST(?::BIGINT[])",
		id, pq.Array(genresIDs),
	).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleOnModerationGenresGenre) {
			return 400, errors.New("жанр не найден")
		}
		return 500, err
	}

	return 0, nil
}
