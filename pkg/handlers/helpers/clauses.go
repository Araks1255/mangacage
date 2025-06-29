package helpers

import "gorm.io/gorm/clause"

var OnConflictClause = clause.OnConflict{
	Columns:   []clause.Column{{Name: "existing_id"}},
	UpdateAll: true,
}
