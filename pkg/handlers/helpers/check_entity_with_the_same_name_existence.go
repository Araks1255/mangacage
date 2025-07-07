package helpers

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func CheckEntityWithTheSameNameExistence(db *gorm.DB, entity string, name, englishName, originalName *string) (existence bool, err error) {
	var dereferencedName, dereferencedEnglishName, dereferencedOriginalName string

	if name != nil {
		dereferencedName = *name
	}
	if englishName != nil {
		dereferencedEnglishName = *englishName
	}
	if originalName != nil {
		dereferencedOriginalName = *originalName
	}

	switch entity {
	case "titles", "volumes", "chapters", "authors":
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE lower(name) = lower(?) OR lower(english_name) = lower(?) OR original_name = ?)", entity)

		err = db.Raw(query, dereferencedName, dereferencedEnglishName, dereferencedOriginalName).Scan(&existence).Error

	case "teams", "genres", "tags":
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE lower(name) = lower(?))", entity)

		err = db.Raw(query, dereferencedName).Scan(&existence).Error

	case "users":
		err = db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE lower(user_name) = lower(?))", dereferencedName).Scan(&existence).Error

	default:
		return true, errors.New("недопустимая сущность")
	}

	if err != nil {
		return true, err
	}

	return existence, nil
}
