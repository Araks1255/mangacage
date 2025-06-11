package titles

import (
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/enum"
)

func ParseInsertError(err error) (code int, reason string) {
	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitlesOnModerationAuthor) {
		return 404, "автор не найден"
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTitlesOnModerationEnglishName) {
		return 409, "тайтл с таким английским названием уже ожидает модерации"
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTitlesOnModerationOriginalName) {
		return 409, "тайтл с таким оригинальным названием уже ожидает модерации"
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniTitlesOnModerationName) {
		return 409, "тайтл с таким названием уже ожидает модерации"
	}

	if dbErrors.IsInvalidEnumValue(err, enum.TitlePublishingStatus) {
		return 409, "указан неверный статус выпуска тайтла"
	}

	if dbErrors.IsInvalidEnumValue(err, enum.TitleTranslatingStatus) {
		return 409, "указан неверный статус перевода тайтла"
	}

	if dbErrors.IsInvalidEnumValue(err, enum.TitleType) {
		return 409, "указан неверный тип тайтла"
	}

	return 500, err.Error()
}
