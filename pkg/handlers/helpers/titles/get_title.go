package titles

import "gorm.io/gorm"

func GetTitle(db *gorm.DB) *gorm.DB {
	return db.Table("titles").Select("*")
}

func GetEditedTitleOnModeration(db *gorm.DB) *gorm.DB {
	return db.Table("titles_on_moderation AS tom").Select("tom.*, t.name AS existing").Joins("INNER JOIN titles AS t ON t.id = tom.existing_id")
}

func GetNewTitleOnModeration(db *gorm.DB) *gorm.DB {
	return db.Table("titles_on_moderation AS tom").Select(
		`tom.*,
		a.name AS author, a.id AS author_id,
		ARRAY_AGG(DISTINCT g.name) AS genres,
		ARRAY_AGG(DISTINCT t.name) AS tags`,
	).
		Joins("INNER JOIN authors AS a ON a.id = tom.author_id").
		Joins("INNER JOIN title_on_moderation_genres AS tomg ON tomg.title_on_moderation_id = tom.id").
		Joins("INNER JOIN genres AS g ON tomg.genre_id = g.id").
		Joins("INNER JOIN title_on_moderation_tags AS tomt ON tomt.title_on_moderation_id = tom.id").
		Joins("INNER JOIN tags AS t ON tomt.tag_id = t.id").
		Group("tom.id, a.id")
}

func GetTitleOnModeration(db *gorm.DB) *gorm.DB {
	return db.Table("titles_on_moderation AS tom").Select(
		`tom.*, titles.name AS existing,
		a.name AS author,
		ARRAY(
			SELECT DISTINCT g.name FROM genres AS g
			INNER JOIN title_on_moderation_genres AS tomg ON tomg.genre_id = g.id
			WHERE tomg.title_on_moderation_id = tom.id
		)::TEXT[] AS genres,
		ARRAY(
			SELECT DISTINCT t.name FROM tags AS t
			INNER JOIN title_on_moderation_tags AS tomt ON tomt.tag_id = t.id
			WHERE tomt.title_on_moderation_id = tom.id
		)::TEXT[] AS tags`, // Тут используется SELECT ARRAY() вместо ARRAY_AGG(), потому что ARRAY_AGG возвращает NULL при отсутствии резульатов, а SELECT ARRAY пустой массив (а NULL не сканируется в поле типа pq.StringArray, даже если это указатель)
	).
		Joins("LEFT JOIN titles ON tom.existing_id = titles.id").
		Joins("LEFT JOIN authors AS a ON tom.author_id = a.id")
}
