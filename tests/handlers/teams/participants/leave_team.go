package participants

import (
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetLeaveTeamScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                     LeaveTeamSuccess(env),
		"by team leader success":      LeaveTeamByTeamLeaderSuccess(env),
		"by only team leader success": LeaveTeamWithOnlyTeamLeaderSuccess(env),
		"with no team":                LeaveTeamWithNoTeam(env),
	}
}

func LeaveTeamSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer", "moder"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, participantID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, participantID, teamID); err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/me", h.LeaveTeam)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/me", nil)

		cookie, err := testhelpers.CreateCookieWithToken(participantID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var participantTeamRoles []string

		if err := env.DB.Raw(
			`SELECT r.name FROM roles AS r
			INNER JOIN user_roles AS ur ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.type = 'team'`,
			participantID,
		).Scan(&participantTeamRoles).Error; err != nil {
			t.Fatal(err)
		}

		if len(participantTeamRoles) != 0 {
			t.Fatal("роли участника не удалились")
		}
	}
}

func LeaveTeamByTeamLeaderSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID, Roles: []string{"typer"}})
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/me", h.LeaveTeam)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/me", nil)

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

		if !slices.Contains(participantRoles, "team_leader") {
			t.Fatal("участник не назначился лидером при уходе лидера")
		}

		if len(participantRoles) != 2 {
			t.Fatal("ещё одна роль участника, назначенного лидером, куда-то делась")
		}
	}
}

func LeaveTeamWithOnlyTeamLeaderSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, leaderID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/me", h.LeaveTeam)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/me", nil)

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

func LeaveTeamByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := participants.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/me", h.LeaveTeam)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/me", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func LeaveTeamWithNoTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		participantID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/my/participants/me", h.LeaveTeam)

		req := httptest.NewRequest("DELETE", "/teams/my/participants/me", nil)

		cookie, err := testhelpers.CreateCookieWithToken(participantID, env.SecretKey)
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
