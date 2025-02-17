package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `gorm:"unique" json:"userName" binding:"required"`
	Password     string `json:"password"`
	AboutYorself string `json:"aboutYourself"`
	Role         string
	TeamID       uint
	team         Team `gorm:"foreignKey:TeamID;references:id"`
}
