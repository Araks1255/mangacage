package testhelpers

import "gorm.io/gorm"

func AddTitleToFavorites(db *gorm.DB, userID, titleID uint) error {
	return db.Exec("INSERT INTO user_favorite_titles (user_id, title_id) VALUES (?, ?)", userID, titleID).Error
}

func AddGenreToFavorites(db *gorm.DB, userID, genreID uint) error {
	return db.Exec("INSERT INTO user_favorite_genres (user_id, genre_id) VALUES (?, ?)", userID, genreID).Error
}

func AddChapterToFavorites(db *gorm.DB, userID, chapterID uint) error {
	return db.Exec("INSERT INTO user_favorite_chapters (user_id, chapter_id) VALUES (?, ?)", userID, chapterID).Error
}
