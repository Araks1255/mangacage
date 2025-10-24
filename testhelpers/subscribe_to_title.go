package testhelpers

import "gorm.io/gorm"

func SubscribeToTitle(db *gorm.DB, titleID, userID uint) error {
	return db.Exec("INSERT INTO user_titles_subscribed_to (user_id, title_id) VALUES (?, ?)", userID, titleID).Error
}
