package participants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetAddRoleToParticipantScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                                  AddRoleToParticipantSuccess(env),
		"success transfer team leader role":        AddRoleOfTeamLeaderToParticipantSuccess(env),
		"unauthorized":                             AddRoleToParticipantByUnauthorizedUser(env),
		"non team leader":                          AddRoleToParticipantByNonTeamLeader(env),
		"invalid participant id":                   AddRoleToParticipantWithInvalidParticipantId(env),
		"invalid role id":                          AddRoleToParticipantWithInvalidRoleId(env),
		"no role id":                               AddRoleToParticipantWithNoRoleId(env),
		"wrong participant id":                     AddRoleToParticipantWithWrongParticipantId(env),
		"wrong role id":                            AddRoleToParticipantWithWrongRoleId(env),
		"role of team leader by ex team leader":    AddRoleOfTeamLeaderToParticipantByExTeamLeader(env),
		"role of ex team leader by ex team leader": AddRoleOfExTeamLeaderByExTeamLeader(env),
		"participant from another team":            AddRoleToParticipantFromAnotherTeam(env),
	}
}

func AddRoleToParticipantSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		var roleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID)
		if roleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleOfTeamLeaderToParticipantSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		var teamLeaderRoleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'team_leader'").Scan(&teamLeaderRoleID)
		if teamLeaderRoleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": teamLeaderRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		req := httptest.NewRequest("POST", "/teams/my/participants/18/roles", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		req := httptest.NewRequest("POST", "/teams/my/participants/18/roles", nil)
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantWithInvalidParticipantId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidParticipantID := "*-*"

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		url := fmt.Sprintf("/teams/my/participants/%s/roles", invalidParticipantID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantWithInvalidRoleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidRoleID := "^_^/"

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]any{
			"roleId": invalidRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/teams/my/participants/18/roles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantWithNoRoleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]any{
			"randomParameter": "?_?",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/teams/my/participants/18/roles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantWithWrongParticipantId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		wrongParticipantID := 9223372036854775807

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": 18,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", wrongParticipantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantWithWrongRoleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": 9223372036854775807,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleOfTeamLeaderToParticipantByExTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		exTeamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"ex_team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, exTeamLeaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, exTeamLeaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		var roleOfTeamLeaderID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'team_leader'").Scan(&roleOfTeamLeaderID)
		if roleOfTeamLeaderID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": roleOfTeamLeaderID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(exTeamLeaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleOfExTeamLeaderByExTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		exTeamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"ex_team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, exTeamLeaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, exTeamLeaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		var roleOfExTeamLeaderID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'ex_team_leader'").Scan(&roleOfExTeamLeaderID)
		if roleOfExTeamLeaderID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": roleOfExTeamLeaderID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(exTeamLeaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddRoleToParticipantFromAnotherTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		anotherTeamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: anotherTeamID})
		if err != nil {
			t.Fatal(err)
		}

		var roleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID)
		if roleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}
