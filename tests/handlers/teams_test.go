package handlers

import (
	"testing"

	"github.com/Araks1255/mangacage/tests/handlers/teams"
	"github.com/Araks1255/mangacage/tests/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/tests/handlers/teams/participants"
)

// Teams
func TestCreateTeam(t *testing.T) {
	scenarios := teams.GetCreateTeamScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestEditTeam(t *testing.T) {
	scenarios := teams.GetEditTeamScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTeamCover(t *testing.T) {
	scenarios := teams.GetGetTeamCoverScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTeam(t *testing.T) {
	scenarios := teams.GetGetTeamScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeleteTeam(t *testing.T) {
	scenarios := teams.GetDeleteTeamScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

// Join requests

func TestAcceptTeamJoinRequest(t *testing.T) {
	scenarios := joinrequests.GetAcceptTeamJoinRequestScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestCancelTeamJoinRequest(t *testing.T) {
	scenarios := joinrequests.GetCancelTeamJoinRequestScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeclineTeamJoinRequest(t *testing.T) {
	scenarios := joinrequests.GetDeclineTeamJoinRequestScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetMyTeamJoinRequests(t *testing.T) {
	scenarios := joinrequests.GetGetMyTeamJoinRequestsScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTeamJoinRequestsOfMyTeam(t *testing.T) {
	scenarios := joinrequests.GetGetTeamJoinRequestsOfMyTeamScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestSubmitTeamJoinRequest(t *testing.T) {
	scenarios := joinrequests.GetSubmitTeamJoinRequestScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

// Participants

func TestAddRoleToParticipant(t *testing.T) {
	scenarios := participants.GetAddRoleToParticipantScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestDeleteParticipantRole(t *testing.T) {
	scenarios := participants.GetDeleteParticipantRoleScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestLeaveTeam(t *testing.T) {
	scenarios := participants.GetLeaveTeamScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestExcludeParticipant(t *testing.T) {
	scenarios := participants.GetExcludeParticipantScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}

func TestGetTeamParticipants(t *testing.T) {
	scenarios := participants.GetGetTeamParticipantsScenarios(env)

	for name, scenario := range scenarios {
		t.Run(name, scenario)
	}
}
