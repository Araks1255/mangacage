package models

type Author struct {
	ID     uint `gorm:"primaryKey;autoIncrement:true"`
	Name   string
	Genres []Genre `gorm:"many2many:author_genres;"`
}
