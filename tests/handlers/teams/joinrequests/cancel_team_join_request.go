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

func GetCancelTeamJoinRequestScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":            CancelTeamJoinRequestSuccess(env),
		"unauthorized":       CancelTeamJoinRequestByUnauthorizedUser(env),
		"from another user":  CancelTeamJoinRequestFromAnotherUser(env),
		"wrong request id":   CancelTeamJoinRequestWithWrongRequestId(env),
		"invalid request id": CancelTeamJoinRequestWithInvalidRequestId(env),
	}
}

func CancelTeamJoinRequestSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		candidateID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, candidateID)
		if err != nil {
			t.Fatal(err)
		}

		requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID, "")
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/join-requests/:id", h.CancelTeamJoinRequest)

		url := fmt.Sprintf("/teams/join-requests/%d", requestID)
		req := httptest.NewRequest("DELETE", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(candidateID, env.SecretKey)

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
	}
}

func CancelTeamJoinRequestByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/join-requests/:id", h.CancelTeamJoinRequest)

		req := httptest.NewRequest("DELETE", "/teams/join-requests/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func CancelTeamJoinRequestFromAnotherUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		candidateID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, candidateID)
		if err != nil {
			t.Fatal(err)
		}

		requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID, "")
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/join-requests/:id", h.CancelTeamJoinRequest)

		url := fmt.Sprintf("/teams/join-requests/%d", requestID)
		req := httptest.NewRequest("DELETE", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
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

func CancelTeamJoinRequestWithWrongRequestId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		requestID := 9223372036854775807

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/join-requests/:id", h.CancelTeamJoinRequest)

		url := fmt.Sprintf("/teams/join-requests/%d", requestID)
		req := httptest.NewRequest("DELETE", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
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

func CancelTeamJoinRequestWithInvalidRequestId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		requestID := "*_*"

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/teams/join-requests/:id", h.CancelTeamJoinRequest)

		url := fmt.Sprintf("/teams/join-requests/%s", requestID)
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
