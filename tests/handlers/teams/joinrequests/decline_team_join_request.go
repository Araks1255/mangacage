package joinrequests

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetDeclineTeamJoinRequestScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                 DeclineTeamJoinRequestSuccess(env),
		"unauthorized":            DeclineTeamJoinRequestByUnauthorizedUser(env),
		"non team leader":         DeclineTeamJoinRequestByNonTeamLeader(env),
		"request to another team": DeclineTeamJoinRequestToAnotherTeam(env),
		"wrong request id":        DeclineTeamJoinRequestWithWrongRequestId(env),
		"invalid request id":      DeclineTeamJoinRequestWithInvalidRequestId(env),
	}
}

func DeclineTeamJoinRequestSuccess(env testenv.Env) func(*testing.T) {
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

		candidateID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID, "")
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "vice_team_leader"}))
		r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%d", requestID)
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
	}
}

func DeclineTeamJoinRequestByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "vice_team_leader"}))
		r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

		req := httptest.NewRequest("DELETE", "/teams/my/join-requests/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeclineTeamJoinRequestByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "vice_team_leader"}))
		r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

		req := httptest.NewRequest("DELETE", "/teams/my/join-requests/18", nil)

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

func DeclineTeamJoinRequestToAnotherTeam(env testenv.Env) func(*testing.T) {
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

		candidateID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		anotherTeamID, err := testhelpers.CreateTeam(env.DB, candidateID)
		if err != nil {
			t.Fatal(err)
		}

		requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, anotherTeamID, "")
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "vice_team_leader"}))
		r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%d", requestID)
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

func DeclineTeamJoinRequestWithWrongRequestId(env testenv.Env) func(*testing.T) {
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

		requestID := 9223372036854775807

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "vice_team_leader"}))
		r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%d", requestID)
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

func DeclineTeamJoinRequestWithInvalidRequestId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidRequestID := ":o"

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "vice_team_leader"}))
		r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%s", invalidRequestID)
		req := httptest.NewRequest("DELETE", url, nil)

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
