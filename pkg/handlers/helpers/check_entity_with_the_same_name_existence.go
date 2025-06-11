package helpers

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func CheckEntityWithTheSameNameExistence(db *gorm.DB, entity, name string, englishName, originalName *string) (existence bool, err error) {
	var dereferencedEnglishName, dereferencedOriginalName string

	if englishName != nil {
		dereferencedEnglishName = *englishName
	}
	if originalName != nil {
		dereferencedOriginalName = *originalName
	}

	switch entity {
	case "titles", "volumes", "chapters":
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE lower(name) = lower(?) OR lower(english_name) = lower(?) OR original_name = ?)", entity)

		err = db.Raw(query, name, dereferencedEnglishName, dereferencedOriginalName).Scan(&existence).Error

	case "teams", "users":
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE lower(name) = lower(?))", entity)

		err = db.Raw(query, name).Scan(&existence).Error

	default:
		return true, errors.New("недопустимая сущность")
	}

	if err != nil {
		return true, err
	}

	return existence, nil
}
