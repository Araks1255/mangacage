package titles

import "gorm.io/gorm"

func GetTitle(db *gorm.DB) *gorm.DB {
	return db.Table("titles").Select("*")
}

func GetTitleWithDependencies(db *gorm.DB) *gorm.DB {
	return db.Table("titles AS t").Select(
		`t.*,
		a.name AS author, a.id AS author_id,
		ARRAY_AGG(DISTINCT g.name)::TEXT[] AS genres,
		ARRAY_AGG(DISTINCT tags.name)::TEXT[] AS tags`,
	).
		Joins("INNER JOIN authors AS a ON t.author_id = a.id").
		Joins("INNER JOIN title_genres AS tg ON tg.title_id = t.id").
		Joins("INNER JOIN genres AS g ON tg.genre_id = g.id").
		Joins("INNER JOIN title_tags AS tt ON t.id = tt.title_id").
		Joins("INNER JOIN tags ON tt.tag_id = tags.id").
		Group("t.id, a.id")
}
