package models

import (
	"time"
)

type TeamJoinRequest struct {
	ID        uint `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time

	IntroductoryMessage *string

	RoleID *uint
	Role   *Role `gorm:"foreignKey:RoleID;references:id;constraint:OnDelete:SET NULL"`

	CandidateID uint  `gorm:"not null"`
	Candidate   *User `gorm:"foreignKey:CandidateID;references:id;constraint:OnDelete:CASCADE"`

	TeamID uint  `gorm:"not null"`
	Team   *Team `gorm:"foreignKey:TeamID;references:id;constraint:OnDelete:CASCADE"`
}
