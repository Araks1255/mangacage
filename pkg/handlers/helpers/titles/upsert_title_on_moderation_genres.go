package titles

import (
	"errors"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func UpsertTitleOnModerationGenres(db *gorm.DB, titleOnModerationID uint, genresIDs []uint) (code int, err error) {
	err = db.Exec("DELETE FROM title_on_moderation_genres WHERE title_on_moderation_id = ?", titleOnModerationID).Error
	if err != nil {
		return 500, err
	}

	err = db.Exec(
		`INSERT INTO title_on_moderation_genres (title_on_moderation_id, genre_id)
		SELECT ?, UNNEST(?::BIGINT[])`,
		titleOnModerationID, pq.Array(genresIDs),
	).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleOnModerationGenresGenre) {
			return 404, errors.New("жанры не найдены")
		} else {
			return 500, err
		}
	}

	return 0, nil
}
