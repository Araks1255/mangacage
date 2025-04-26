package models

import "github.com/lib/pq"

type Author struct {
	ID     uint   `gorm:"primaryKey;autoIncrement:true"`
	Name   string `gorm:"unique"`
	About  string
	Genres []Genre `gorm:"many2many:author_genres;constraint:OnDelete:CASCADE"`
}

type AuthorDTO struct {
	ID     uint           `json:"id"`
	Name   string         `json:"name"`
	About  string         `json:"about,omitempty"`
	Genres pq.StringArray `json:"genres,omitempty" gorm:"type:TEXT[]"`
}
