package models

import "time"

type UserViewedChapter struct {
	CreatedAt time.Time

	UserID uint
	User   *User `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`

	ChapterID uint
	Chapter   *Chapter `gorm:"foreignKey:ChapterID;references:id;constraint:OnDelete:CASCADE"`
}
