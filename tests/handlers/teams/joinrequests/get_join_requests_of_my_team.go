package joinrequests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTeamJoinRequestsOfMyTeamScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      GetTeamJoinRequestsOfMyTeamSuccess(env),
		"unauthorized": GetTeamJoinRequestsOfMyTeamByUnauthorizedUser(env),
		"not in team":  GetTeamJoinRequestsOfMyTeamByUserThatNotInTeam(env),
		"no requests":  GetTeamJoinRequestsOfMyTeamWithNoRequests(env),
	}
}

func GetTeamJoinRequestsOfMyTeamSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			candidateID, err := testhelpers.CreateUser(env.DB)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID, "typer"); err != nil {
				t.Fatal(err)
			}
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/my/join-requests", h.GetTeamJoinRequestsOfMyTeam)

		req := httptest.NewRequest("GET", "/teams/my/join-requests", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) < 2 {
			t.Fatal("дошли не все запросы")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id запроса не дошел")
		}
		if _, ok := resp[0]["roleId"]; !ok {
			t.Fatal("id роли не дошло")
		}
		if role, ok := resp[0]["role"]; !ok || role != "typer" {
			t.Fatal("роль не дошла (или вернулась неверная)")
		}
		if _, ok := resp[0]["candidateId"]; !ok {
			t.Fatal("id кандидата не вернулось")
		}
		if _, ok := resp[0]["candidate"]; !ok {
			t.Fatal("имя кандидата не вернулось")
		}
	}
}

func GetTeamJoinRequestsOfMyTeamByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/my/join-requests", h.GetTeamJoinRequestsOfMyTeam)

		req := httptest.NewRequest("GET", "/teams/my/join-requests", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTeamJoinRequestsOfMyTeamByUserThatNotInTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/my/join-requests", h.GetTeamJoinRequestsOfMyTeam)

		req := httptest.NewRequest("GET", "/teams/my/join-requests", nil)

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

func GetTeamJoinRequestsOfMyTeamWithNoRequests(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/my/join-requests", h.GetTeamJoinRequestsOfMyTeam)

		req := httptest.NewRequest("GET", "/teams/my/join-requests", nil)

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
