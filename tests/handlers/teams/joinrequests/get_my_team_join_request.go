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

func GetGetMyTeamJoinRequestsScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      GetMyTeamJoinRequestsSuccess(env),
		"unauthorized": GetMyTeamJoinRequestsByUnauthorizedUser(env),
		"no requests":  GetMyTeamJoinRequestsWithNoJoinRequests(env),
	}
}

func GetMyTeamJoinRequestsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			teamID, err := testhelpers.CreateTeam(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := testhelpers.CreateTeamJoinRequest(env.DB, userID, teamID, "typer"); err != nil {
				t.Fatal(err)
			}
		}

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/join-requests/my", h.GetMyTeamJoinRequests)

		req := httptest.NewRequest("GET", "/teams/join-requests/my", nil)

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
			t.Fatal("вернулись не все запросы")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не вернулся")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не вернулось")
		}
		if role, ok := resp[0]["role"]; !ok || role != "typer" {
			t.Fatal("роль не вернулась (или вернулась неверная)")
		}
		if _, ok := resp[0]["roleId"]; !ok {
			t.Fatal("id роли не вернулось")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("название команды не вернулось")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не вернулось")
		}
	}
}

func GetMyTeamJoinRequestsByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/join-requests/my", h.GetMyTeamJoinRequests)

		req := httptest.NewRequest("GET", "/teams/join-requests/my", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetMyTeamJoinRequestsWithNoJoinRequests(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/teams/join-requests/my", h.GetMyTeamJoinRequests)

		req := httptest.NewRequest("GET", "/teams/join-requests/my", nil)

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
