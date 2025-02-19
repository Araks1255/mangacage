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

	TeamID sql.NullInt32
	Team   Team `gorm:"foreignKey:TeamID;references:id"`

	Roles []Role `gorm:"many2many:user_roles;"`
}
