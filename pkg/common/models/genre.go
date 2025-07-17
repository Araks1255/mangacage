package models

type Genre struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"not null"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type GenreOnModeration struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"not null"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraints:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (GenreOnModeration) TableName() string {
	return "genres_on_moderation"
}

func (g GenreOnModeration) ToGenre() Genre {
	return Genre{
		Name:        g.Name,
		ModeratorID: g.ModeratorID,
	}
}
