package titles

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTitleTranslateRequests(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      GetTitleTranslateRequestsSuccess(env),
		"without team": GetTitleTranslateRequestsWithoutTeam(env),
	}
}

func GetTitleTranslateRequestsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := testhelpers.CreateTitleTranslateRequest(env.DB, titleID, teamID, "message"); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/titles/translate-requests", h.GetTitleTranslateRequests)

		req := httptest.NewRequest("GET", "/titles/translate-requests", nil)

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

		if len(resp) != 2 {
			t.Fatal("не все запросы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время отправки не дошло")
		}
		if _, ok := resp[0]["message"]; !ok {
			t.Fatal("сообщение не дошло")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
	}
}

func GetTitleTranslateRequestsWithoutTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/titles/translate-requests", h.GetTitleTranslateRequests)

		req := httptest.NewRequest("GET", "/titles/translate-requests", nil)

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
