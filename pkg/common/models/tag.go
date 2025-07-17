package models

type Tag struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"not null"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

type TagOnModeration struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Name string `gorm:"not null"`

	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraints:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`
}

func (TagOnModeration) TableName() string {
	return "tags_on_moderation"
}

func (t TagOnModeration) ToTag() Tag {
	return Tag{
		Name:        t.Name,
		ModeratorID: t.ModeratorID,
	}
}
