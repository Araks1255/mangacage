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

func GetDeleteParticipantRoleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                               DeleteParticipantRoleSuccess(env),
		"self role success":                     DeleteSelfRoleSuccess(env),
		"unauthorized":                          DeleteParticipantRoleByUnauthorizedUser(env),
		"non team leader":                       DeleteParticipantRoleByNonTeamLeader(env),
		"self ex team leader role":              DeleteSelfExTeamLeaderRoleByExTeamLeader(env),
		"invalid participant id":                DeleteParticipantRoleWithInvalidParticipantId(env),
		"invalid role id":                       DeleteParticipantRoleWithInvalidRoleId(env),
		"wrong participant id":                  DeleteParticipantRoleWithWrongParticipantId(env),
		"wrong role id":                         DeleteParticipantRoleWithWrongRoleId(env),
		"role does not exist":                   DeleteParticipantRoleThatDoesNotExist(env),
		"self team leader role":                 DeleteSelfTeamLeaderRole(env),
		"team leader role by ex team leader":    DeleteTeamLeaderRoleByExTeamLeader(env),
		"participant from another team":         DeleteParticipantFromAnotherTeamRole(env),
		"ex team leader role by ex team leader": DeleteParticipantExTeamLeaderRoleByExTeamLeader(env),
	}
}

func DeleteParticipantRoleSuccess(env testenv.Env) func(*testing.T) {
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

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}, TeamID: teamID})
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
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteSelfRoleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader", "typer"}})
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

		var roleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID)
		if roleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", leaderID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteParticipantRoleByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/18/roles", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteParticipantRoleByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/18/roles", nil)

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

func DeleteSelfRoleByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, participantID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, participantID, teamID); err != nil {
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
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(participantID, env.SecretKey)
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

func DeleteSelfExTeamLeaderRoleByExTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		exLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"ex_team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, exLeaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, exLeaderID, teamID); err != nil {
			t.Fatal(err)
		}

		var exTeamLeaderRoleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'ex_team_leader'").Scan(&exTeamLeaderRoleID)
		if exTeamLeaderRoleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": exTeamLeaderRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", exLeaderID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(exLeaderID, env.SecretKey)
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

func DeleteParticipantRoleWithInvalidParticipantId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidParticipantID := "X_X"

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		url := fmt.Sprintf("/teams/my/participants/%s/roles", invalidParticipantID)
		req := httptest.NewRequest("DELETE", url, nil)
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

func DeleteParticipantRoleWithInvalidRoleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidRoleID := "I_I"

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]any{
			"roleId": invalidRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("DELETE", "/teams/my/participants/18/roles", bytes.NewBuffer(jsonBody))
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

func DeleteParticipantRoleWithWrongParticipantId(env testenv.Env) func(*testing.T) {
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

		participantID := 9223372036854775807

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": 18,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
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

func DeleteParticipantRoleWithWrongRoleId(env testenv.Env) func(*testing.T) {
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

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}, TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		roleID := 9223372036854775807

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]int{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
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

func DeleteParticipantRoleThatDoesNotExist(env testenv.Env) func(*testing.T) {
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
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 409 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteParticipantExTeamLeaderRoleByExTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"ex_team_leader"}})
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

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"ex_team_leader"}, TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		var exTeamLeaderRoleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'ex_team_leader'").Scan(&exTeamLeaderRoleID)
		if exTeamLeaderRoleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": exTeamLeaderRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
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

func DeleteSelfTeamLeaderRole(env testenv.Env) func(*testing.T) {
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

		var teamLeaderRoleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'team_leader'").Scan(&teamLeaderRoleID)
		if teamLeaderRoleID == 0 {
			t.Fatal("не удалось получить роль")
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": teamLeaderRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", leaderID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
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

func DeleteTeamLeaderRoleByExTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		exLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"ex_team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, exLeaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, exLeaderID, teamID); err != nil {
			t.Fatal(err)
		}

		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}, TeamID: teamID})
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
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": teamLeaderRoleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", leaderID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(exLeaderID, env.SecretKey)
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

func DeleteParticipantFromAnotherTeamRole(env testenv.Env) func(*testing.T) {
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

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}, TeamID: anotherTeamID})
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
		r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

		body := map[string]uint{
			"roleId": roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
		req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
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
