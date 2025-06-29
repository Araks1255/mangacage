package models

type TitleRate struct {
	TitleID uint
	Title Title `gorm:"not null;foreignKey:TitleID;references:id;constraint:OnDelete:CASCADE"`

	UserID uint
	User User `gorm:"not null;foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`

	Rate int `gorm:"not null;check:rate >= 1 AND rate <= 5"`
}