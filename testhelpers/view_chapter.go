package testhelpers

import "gorm.io/gorm"

func ViewChapter(db *gorm.DB, viewerID, chapterID uint) error { // Под "посмотреть главу" подразумевается добавить просмотр от пользователя к главе. На этих просмотрах работает история чтения и сбор статистики по самым популярным тайтлам
	if result := db.Exec(
		"INSERT INTO user_viewed_chapters (user_id, chapter_id) VALUES (?, ?)",
		viewerID, chapterID,
	); result.Error != nil {
		return result.Error
	}
	return nil
}
