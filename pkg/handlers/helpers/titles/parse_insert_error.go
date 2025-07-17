package titles

import (
	"errors"

	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/enum"
)

func ParseTitleOnModerationInsertError(err error) (code int, formatedErr error) {
	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitlesOnModerationAuthor) {
		return 404, errors.New("автор не найден")
	}

	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitlesOnModerationAuthorOnModeration) {
		return 404, errors.New("автор на модерации не найден")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTitleOnModerationEnglishName) {
		return 409, errors.New("тайтл с таким английским названием уже ожидает модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTitleOnModerationOriginalName) {
		return 409, errors.New("тайтл с таким оригинальным названием уже ожидает модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTitleOnModerationName) {
		return 409, errors.New("тайтл с таким названием уже ожидает модерации")
	}

	if dbErrors.IsInvalidEnumValue(err, enum.TitlePublishingStatus) {
		return 409, errors.New("указан неверный статус выпуска тайтла")
	}

	if dbErrors.IsInvalidEnumValue(err, enum.TitleTranslatingStatus) {
		return 409, errors.New("указан неверный статус перевода тайтла")
	}

	if dbErrors.IsInvalidEnumValue(err, enum.TitleType) {
		return 409, errors.New("указан неверный тип тайтла")
	}

	return 500, err
}

func ParseTitleTeamInsertError(err error) (code int, formatedErr error) {
	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleTeamsTitle) {
		return 404, errors.New("тайтл не найден")
	}

	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitlesTeamsTeam) {
		return 409, errors.New("вы не состоите в команде перевода")
	}

	if dbErrors.IsUniqueViolation(err, constraints.TitleTeamsPkey) {
		return 409, errors.New("ваша команда уже переводит этот тайтл")
	}

	return 500, err
}

func ParseTitleTranslateRequestInsertError(err error) (code int, formatedErr error) {
	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleTranslateRequestTitle) {
		return 404, errors.New("тайтл не найден")
	}

	if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleTranslateRequestTeam) {
		return 404, errors.New("команда не найдена") // мало ли
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTitleTranslateRequest) {
		return 409, errors.New("ваша команда уже подала заявку на перевод этого тайтла")
	}

	return 500, err
}
