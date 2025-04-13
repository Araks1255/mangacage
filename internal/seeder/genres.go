package seeder

import "gorm.io/gorm"

func seedGenres(db *gorm.DB) error {
	if result := db.Exec(
		`INSERT INTO genres (name) VALUES
		('fighting'), ('action'), ('romance'), ('dystopia')
		ON CONFLICT DO NOTHING`,
	); result.Error != nil {
		return result.Error
	}
	return nil
}