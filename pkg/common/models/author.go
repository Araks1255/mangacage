package models

type Author struct {
	ID     uint   `gorm:"primaryKey;autoIncrement:true"`
	Name   string `gorm:"unique"`
	About  string
	Genres []Genre `gorm:"many2many:author_genres;constraint:OnDelete:CASCADE"`
}
