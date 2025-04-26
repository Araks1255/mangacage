package testhelpers

import "gorm.io/gorm"

func TranslateTitle(db *gorm.DB, teamID, titleID uint) error {
	if result := db.Exec("UPDATE titles SET team_id = ? WHERE id = ?", teamID, titleID); result.Error != nil {
		return result.Error
	}
	return nil
}