package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `gorm:"unique"`
	Password     string
	AboutYorself string
	TgUserID     int64

	TeamID sql.NullInt64
	Team   Team `gorm:"foreignKey:TeamID;references:id;OnDelete:SET NULL"`

	Roles                    []Role  `gorm:"many2many:user_roles;OnDelete:SET NULL"`
	TitlesUserIsSubscribedTo []Title `gorm:"many2many:user_titles_subscribed_to;OnDelete:SET NULL"`
}
