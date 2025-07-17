package dto

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
)

type CreateTeamJoinRequestDTO struct {
	IntroductoryMessage *string `json:"introductoryMessage"`
	RoleID              *uint   `json:"roleId"`
}

type ResponseTeamJoinRequestDTO struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	IntroductoryMessage *string `json:"introductoryMessage,omitempty"`

	Role   *string `json:"role,omitempty"`
	RoleID *uint   `json:"roleId,omitempty"`

	Candidate   *string `json:"candidate,omitempty"`
	CandidateID *uint   `json:"candidateId,omitempty"`

	Team   *string `json:"team,omitempty"`
	TeamID *uint   `json:"teamId,omitempty"`
}

func (tjr CreateTeamJoinRequestDTO) ToTeamJoinRequest(candidateID, teamID uint) models.TeamJoinRequest {
	return models.TeamJoinRequest{
		IntroductoryMessage: tjr.IntroductoryMessage,
		RoleID:              tjr.RoleID,
		TeamID:              teamID,
		CandidateID:         candidateID,
	}
}
