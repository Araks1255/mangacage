package helpers

import "gorm.io/gorm/clause"

var OnExistingIDConflictClause = clause.OnConflict{
	Columns:   []clause.Column{{Name: "existing_id"}},
	UpdateAll: true,
}

var OnIDConflictClause = clause.OnConflict{
	Columns:   []clause.Column{{Name: "id"}},
	UpdateAll: true,
}
