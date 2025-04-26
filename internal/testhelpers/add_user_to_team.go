package testhelpers

import "gorm.io/gorm"

func AddUserToTeam(db *gorm.DB, userID, teamID uint) error {
	if result := db.Exec("UPDATE users SET team_id = ? WHERE id = ?", teamID, userID); result.Error != nil {
		return result.Error
	}
	return nil
}
