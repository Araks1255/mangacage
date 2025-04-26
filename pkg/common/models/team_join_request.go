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

	CandidateID uint
	User        *User `gorm:"foreignKey:CandidateID;references:id;OnDelete:CASCADE"`

	TeamID uint
	Team   *Team `gorm:"foreignKey:TeamID;references:id;OnDelete:CASCADE"`
}

type TeamJoinRequestDTO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`

	Role   string `json:"role"`
	RoleID uint   `json:"roleId"`

	Candidate   string `json:"candidate"`
	CandidateID uint   `json:"candidateId"`

	Team   string `json:"team"`
	TeamID uint   `json:"teamId"`
}
