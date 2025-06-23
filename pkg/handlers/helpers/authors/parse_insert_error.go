package authors

import (
	"errors"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
)

func ParseAuthorOnModerationInsertError(err error) (code int, formatedErr error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniqAuthorOnModerationName) {
		return 409, errors.New("автор с таким именем уже ожидает модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqAuthorOnModerationEnglishName) {
		return 409, errors.New("автор с таким английским именем уже ожидает модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniAuthorOnModerationOriginalName) {
		return 409, errors.New("автор с таким оригинальным именем уже ожидает модерации")
	}

	return 500, err
}
