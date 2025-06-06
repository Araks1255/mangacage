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

func GetAcceptTeamJoinRequestScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":            AcceptTeamJoinRequestSuccess(env),
		"unauthorized":       AcceptTeamJoinRequestByUnauthorizedUser(env),
		"non team leader":    AcceptTeamJoinRequestByNonTeamLeader(env),
		"to another team":    AcceptTeamJoinRequestToAnotherTeam(env),
		"wrong request id":   AcceptTeamJoinRequestWithWrongRequestId(env),
		"invalid request id": AcceptTeamJoinRequestWithInvalidRequestId(env),
	}
}

func AcceptTeamJoinRequestSuccess(env testenv.Env) func(*testing.T) {
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

		requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID, "typer")
		if err != nil {
			t.Fatal(err)
		}

		anotherTeamID, err := testhelpers.CreateTeam(env.DB, candidateID)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, anotherTeamID, ""); err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%d/accept", requestID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
	}
}

func AcceptTeamJoinRequestByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

		req := httptest.NewRequest("POST", "/teams/my/join-requests/18/accept", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func AcceptTeamJoinRequestByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

		req := httptest.NewRequest("POST", "/teams/my/join-requests/18/accept", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatal(w.Body.String())
		}
	}
}

func AcceptTeamJoinRequestToAnotherTeam(env testenv.Env) func(*testing.T) {
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

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%d/accept", requestID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func AcceptTeamJoinRequestWithWrongRequestId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		requestID := 9223372036854775807

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%d/accept", requestID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func AcceptTeamJoinRequestWithInvalidRequestId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidRequestID := "$_$"

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

		url := fmt.Sprintf("/teams/my/join-requests/%s/accept", invalidRequestID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
