package seeder

import "gorm.io/gorm"

func seedVolumes(db *gorm.DB) error {
	if result := db.Exec(
		`INSERT INTO volumes (name, description, title_id, creator_id, moderator_id)
		SELECT 'volume_test', '', t.id, u.id, u.id FROM
		(SELECT id FROM titles WHERE name = 'title_test') AS t,
		(SELECT id FROM users WHERE user_name = 'user_test') AS u
		ON CONFLICT DO NOTHING`,
	); result.Error != nil {
		return result.Error
	}
	return nil
}
