package seeder

import "gorm.io/gorm"

func seedAuthors(db *gorm.DB) error {
	if result := db.Exec(
		`INSERT INTO authors (name)
		VALUES ('author_test')
		ON CONFLICT DO NOTHING`,
	); result.Error != nil {
		return result.Error
	}
	return nil
}
