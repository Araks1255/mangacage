package titles

import (
	"errors"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func UpsertTitleOnModerationTags(db *gorm.DB, id uint, tagsIDs []uint) (code int, err error) {
	if len(tagsIDs) == 0 {
		return 0, nil
	}

	query :=
		`WITH delete_tags AS (
			DELETE FROM
				title_on_moderation_tags
			WHERE
				title_on_moderation_id = ?
		)
		INSERT INTO
			title_on_moderation_tags (title_on_moderation_id, tag_id)
		SELECT
			?, UNNEST(?::BIGINT[])`

	if err := db.Exec(query, id, id, pq.Array(tagsIDs)).Error; err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleOnModerationTagsTag) {
			return 400, errors.New("тег не найден")
		}
		return 500, err
	}

	return 0, nil
}
