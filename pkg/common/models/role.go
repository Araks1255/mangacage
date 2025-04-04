package models

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"unique"`
	Type string // Роли на сайте и роли в команде будут разделены, это нужно для удобной выборки ролей только в рамках команды
}
