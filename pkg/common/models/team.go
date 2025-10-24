package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model

	Name                 string `gorm:"not null"`
	Description          string
	NumberOfParticipants uint `gorm:"not null;default:0"`

	CreatorID *uint
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:SET NULL"`

	ModeratorID *uint
	Moderator   *User `gorm:"foreignKey:ModeratorID;references:id;constraint:OnDelete:SET NULL"`

	EditorID *uint
	Editor   *User `gorm:"foreignKey:EditorID;references:id;constraint:OnDelete:SET NULL"`

	CoverPath string `gorm:"not null"`
}

type TeamOnModeration struct {
	gorm.Model
	Name        *string
	Description *string

	ExistingID *uint `gorm:"unique"`
	Team       *Team `gorm:"foreignKey:ExistingID;references:id;constraint:OnDelete:CASCADE"`

	CreatorID uint  `gorm:"unique;not null"`
	Creator   *User `gorm:"foreignKey:CreatorID;references:id;constraint:OnDelete:CASCADE"`

	ModeratorID *uint

	CoverPath *string
}

func (TeamOnModeration) TableName() string {
	return "teams_on_moderation"
}

func (t TeamOnModeration) ToTeam() *Team {
	team := &Team{ModeratorID: t.ModeratorID}

	if t.Name != nil {
		team.Name = *t.Name
	}
	if t.Description != nil {
		team.Description = *t.Description
	}

	if t.ExistingID == nil {
		team.CreatorID = &t.CreatorID
	} else {
		team.EditorID = &t.CreatorID
	}

	return team
}

func (t *TeamOnModeration) SetID(id uint) {
	t.ID = id
}
