package titles

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/titles/translaterequests"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetQuitTranslatingTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                              QuitTranslatingTitleSuccess(env),
		"success by the only one team success": QuitTranslatingTitleByTheOnlyOneTeamSuccess(env),
		"title is not translating":             QuitTranslatingTitleThatIsNotTranslating(env),
	}
}

func QuitTranslatingTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		{
			anotherTeamID, err := testhelpers.CreateTeam(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := testhelpers.TranslateTitle(env.DB, anotherTeamID, titleID); err != nil {
				t.Fatal(err)
			}
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/titles/:id/quit-translating", h.QuitTranslatingTitle)

		url := fmt.Sprintf("/titles/%d/quit-translating", titleID)
		req := httptest.NewRequest("DELETE", url, nil)

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
	}
}

func QuitTranslatingTitleByTheOnlyOneTeamSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/titles/:id/quit-translating", h.QuitTranslatingTitle)

		url := fmt.Sprintf("/titles/%d/quit-translating", titleID)
		req := httptest.NewRequest("DELETE", url, nil)

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

		var status string
		if err := env.DB.Raw("SELECT translating_status FROM titles WHERE id = ?", titleID).Scan(&status).Error; err != nil {
			t.Fatal(err)
		}

		if status != "free" {
			t.Fatal("статус перевода не изменился")
		}
	}
}

func QuitTranslatingTitleThatIsNotTranslating(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/titles/:id/quit-translating", h.QuitTranslatingTitle)

		url := fmt.Sprintf("/titles/%d/quit-translating", titleID)
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
