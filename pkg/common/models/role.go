package models

type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"unique"`
	Type string `gorm:"type:role_type"`
}
