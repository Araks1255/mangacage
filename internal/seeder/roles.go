package seeder

import "gorm.io/gorm"

func seedRoles(db *gorm.DB) error {
	if result := db.Exec(
		`INSERT INTO roles (name, type) VALUES
		('user', 'site'), ('admin', 'site'), ('moder', 'site'),
		('team_leader', 'team'), ('ex_team_leader', 'team'),
		('translater', 'team'), ('typer', 'team'), ('cleaner', 'team')
		ON CONFLICT DO NOTHING`,
	); result.Error != nil {
		return result.Error
	}
	return nil
}
