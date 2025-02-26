package models

type Genre struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"unique"`
}
