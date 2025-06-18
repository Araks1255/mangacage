package titles

import "gorm.io/gorm"

func IsUserTeamTranslatingTitle(db *gorm.DB, userID, titleID uint) (bool, error) {
	var res bool

	err := db.Raw(
		"SELECT EXISTS(SELECT 1 FROM title_teams WHERE title_id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?))",
		titleID, userID,
	).Scan(&res).Error

	if err != nil {
		return false, err
	}

	return res, nil
}
