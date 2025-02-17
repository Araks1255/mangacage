package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `gorm:"unique"`
	AboutYorself string
	Role         string
	TeamID       uint
	team         Team `gorm:"foreignKey:TeamID;references:id"`
}
