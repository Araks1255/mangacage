package models

type Tag struct {
	ID   uint `gorm:"primaryKey;autoIncrement:true"`
	Name string
}

type TagDTO struct {
	Name string `json:"name" form:"name"`
}
