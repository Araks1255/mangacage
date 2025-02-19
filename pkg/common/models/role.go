package models

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"unique"`
}
