package testhelpers

import "gorm.io/gorm"

func TranslateTitle(db *gorm.DB, teamID, titleID uint) error {
	return db.Exec("INSERT INTO title_teams (title_id, team_id) VALUES (?, ?)", titleID, teamID).Error
}
