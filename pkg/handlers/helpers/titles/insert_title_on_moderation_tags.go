package titles

import (
	"errors"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func InsertTitleOnModerationTags(db *gorm.DB, titleOnModerationID uint, tagsIDs []uint) (code int, err error) {
	err = db.Exec(
		`INSERT INTO title_on_moderation_tags (title_on_moderation_id, tag_id)
		SELECT ?, UNNEST(?::BIGINT[])`,
		titleOnModerationID, pq.Array(tagsIDs),
	).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleOnModerationTagsTag) {
			return 404, errors.New("теги не найдены")
		} else {
			return 500, err
		}
	}

	return 0, nil
}
