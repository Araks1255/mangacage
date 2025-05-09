package models

import (
	"database/sql"
	"time"
)

type TeamJoinRequest struct {
	ID                  uint `gorm:"primaryKey;autoIncrement:true"`
	CreatedAt           time.Time
	IntroductoryMessage string

	RoleID sql.NullInt64
	Role   *Role `gorm:"foreignKey:RoleID;references:id;OnDelete:SET NULL"`

	CandidateID uint  `gorm:"NOT NULL"`
	User        *User `gorm:"foreignKey:CandidateID;references:id;OnDelete:CASCADE"`

	TeamID uint  `gorm:"NOT NULL"`
	Team   *Team `gorm:"foreignKey:TeamID;references:id;OnDelete:CASCADE"`
}

type TeamJoinRequestDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	IntroductoryMessage string `json:"introductoryMessage,omitempty"`

	Role   string `json:"role,omitempty"`
	RoleID uint   `json:"roleId,omitempty"`

	Candidate   string `json:"candidate,omitempty"`
	CandidateID uint   `json:"candidateId,omitempty"`

	Team   string `json:"team,omitempty"`
	TeamID uint   `json:"teamId,omitempty"`
}
