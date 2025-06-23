package models

type Tag struct {
	ID   uint `gorm:"primaryKey;autoIncrement:true"`
	Name string
}

type TagDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TagOnModeration struct {
	ID        uint `gorm:"primaryKey;autoIncrement:true"`
	Name      string
	CreatorID uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraints:OnDelete:CASCADE"`
}

func (TagOnModeration) TableName() string {
	return "tags_on_moderation"
}

type TagOnModerationDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name" binding:"required"`
	CreatorID *uint  `json:"creatorId"`
}

func (t TagOnModerationDTO) ToTagOnModeration(creatorID uint) TagOnModeration {
	return TagOnModeration{
		Name:      t.Name,
		CreatorID: creatorID,
	}
}
