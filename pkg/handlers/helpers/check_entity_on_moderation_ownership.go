package helpers

import (
	"fmt"

	"gorm.io/gorm"
)

func CheckEntityOnModerationOwnership(db *gorm.DB, entity string, entityOnModerationID, userID uint) (bool, error) {
	var (
		res bool
		err error
	)

	switch entity {
	case "titles", "chapters", "teams":
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s_on_moderation WHERE id = ? AND creator_id = ?)", entity)
		err = db.Raw(query, entityOnModerationID, userID).Scan(&res).Error

	case "users":
		err = db.Raw(
			"SELECT EXISTS(SELECT 1 FROM users_on_moderation WHERE id = ? AND existing_id = ?)",
			entityOnModerationID, userID,
		).Scan(&res).Error
	}

	if err != nil {
		return false, err
	}

	return res, nil
}
