package participants

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetExcludeParticipantScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":           ExcludeParticipantSuccess(env),
		"invalid id":        ExcludeParticipantWithInvalidId(env),
		"wrong id":          ExcludeParticipantWithWrongId(env),
		"from another team": ExcludeParticipantFromAnotherTeam(env),
		"yourself":          ExcludeYourself(env),
	}
}

func ExcludeParticipantSuccess(env testenv.Env) func(*testing.T) {
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

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer", "translater"}, TeamID: teamID})
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my/participants/:id", h.ExcludeParticipant)

		url := fmt.Sprintf("/teams/my/participants/%d", participantID)
		req := httptest.NewRequest("DELETE", url, nil)

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

		var participantRoles []string

		if err := env.DB.Raw(
			`SELECT r.name FROM roles AS r
			INNER JOIN user_roles AS ur ON ur.role_id = r.id
			WHERE ur.user_id = ?`, participantID,
		).Scan(&participantRoles).Error; err != nil {
			t.Fatal(err)
		}

		if len(participantRoles) != 0 {
			t.Fatal("роли участника не удалились")
		}
	}
}

func ExcludeParticipantWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidParticipantID := "}._.{"

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my/participants/:id", h.ExcludeParticipant)

		url := fmt.Sprintf("/teams/my/participants/%s", invalidParticipantID)
		req := httptest.NewRequest("DELETE", url, nil)

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

func ExcludeParticipantWithWrongId(env testenv.Env) func(*testing.T) {
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
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my/participants/:id", h.ExcludeParticipant)

		url := fmt.Sprintf("/teams/my/participants/%d", participantID)
		req := httptest.NewRequest("DELETE", url, nil)

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

func ExcludeParticipantFromAnotherTeam(env testenv.Env) func(*testing.T) {
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

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my/participants/:id", h.ExcludeParticipant)

		url := fmt.Sprintf("/teams/my/participants/%d", participantID)
		req := httptest.NewRequest("DELETE", url, nil)

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

func ExcludeYourself(env testenv.Env) func(*testing.T) {
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

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my/participants/:id", h.ExcludeParticipant)

		url := fmt.Sprintf("/teams/my/participants/%d", leaderID)
		req := httptest.NewRequest("DELETE", url, nil)

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
