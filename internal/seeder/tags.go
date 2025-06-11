package seeder

import "gorm.io/gorm"

func seedTags(db *gorm.DB) error {
	if result := db.Exec(
		`INSERT INTO tags (name) VALUES
		('maids'), ('loli'), ('reincarnation'), ('japan')
		ON CONFLICT DO NOTHING`,
	); result.Error != nil {
		return result.Error
	}
	return nil
}
