package models

type Genre struct {
	ID   uint `gorm:"primaryKey;autoIncrement:true"`
	Name string
}

type GenreDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GenreOnModeration struct {
	ID   uint `gorm:"primaryKey;autoIncrement:true"`
	Name string

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraints:OnDelete:CASCADE"`
}

func (GenreOnModeration) TableName() string {
	return "genres_on_moderation"
}

type GenreOnModerationDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" binding:"required"`
	CreatorID uint   `json:"creatorId"`
}

func (g GenreOnModerationDTO) ToGenreOnModeration(creatorID uint) GenreOnModeration {
	return GenreOnModeration{
		Name:      g.Name,
		CreatorID: creatorID,
	}
}
