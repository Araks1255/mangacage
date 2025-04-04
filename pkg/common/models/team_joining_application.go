package models

import "time"

type TeamJoiningApplication struct {
	ID                  uint `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt           time.Time
	Role                string
	IntroductoryMessage string

	CandidateID uint
	User        *User `gorm:"foreignKey:CandidateID;references:id;OnDelete:CASCADE"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;OnDelete:CASCADE"`
}
