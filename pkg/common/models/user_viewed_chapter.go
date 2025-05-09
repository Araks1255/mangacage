package models

import "time"

type UserViewedChapter struct {
	CreatedAt time.Time
	UserID    uint
	User      *User `gorm:"foreignKey:UserID;references:id;OnDelete:CASCADE"` // Над OnDelete подумать
	ChapterID uint
	Chapter   *Chapter `gorm:"foreignKey:ChapterID;references:id;OnDelete:CASCADE"`
}
