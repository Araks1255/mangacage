package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string `gorm:"unique" json:"userName" binding:"required"`
	Password     string `json:"password"`
	AboutYorself string `json:"aboutYourself"`
	Role         string
	TeamID       sql.NullInt32
	Team         Team `gorm:"foreignKey:TeamID;references:id"`
}
