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

func GetTranslateTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success free":         TranslateFreeTitleSuccess(env),
		"success request":      TranslateTitleWithTranslateRequestSuccess(env),
		"already translating":  TranslateAlreadyTranslatingTitle(env),
		"request already sent": TranslateTitleWithAlreadySentRequest(env),
		"twice":                TranslateTitleTwice(env),
	}
}

func TranslateFreeTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
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

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/translate", h.TranslateTitle)

		url := fmt.Sprintf("/titles/%d/translate", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}

		var status string
		if err := env.DB.Raw("SELECT translating_status FROM titles WHERE id = ?", titleID).Scan(&status).Error; err != nil {
			t.Fatal(err)
		}

		if status != "ongoing" {
			t.Fatal("статус перевода не изменился")
		}
	}
}

func TranslateTitleWithTranslateRequestSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
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

		anotherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, anotherUserID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/translate", h.TranslateTitle)

		url := fmt.Sprintf("/titles/%d/translate", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}

		var check bool
		if err := env.DB.Raw(
			"SELECT EXISTS(SELECT 1 FROM titles_translate_requests WHERE title_id = ? AND team_id = ?)",
			titleID, teamID,
		).Scan(&check).Error; err != nil {
			t.Fatal(err)
		}

		if !check {
			t.Fatal("заявка на перевод не оставилась")
		}
	}
}

func TranslateAlreadyTranslatingTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
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

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.TranslateTitle(env.DB, teamID, titleID); err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/translate", h.TranslateTitle)

		url := fmt.Sprintf("/titles/%d/translate", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
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

func TranslateTitleWithAlreadySentRequest(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
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

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := testhelpers.CreateTitleTranslateRequest(env.DB, titleID, teamID); err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/translate", h.TranslateTitle)

		url := fmt.Sprintf("/titles/%d/translate", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
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

func TranslateTitleTwice(env testenv.Env) func(*testing.T) {
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
		r.POST("/titles/:id/translate", h.TranslateTitle)

		url := fmt.Sprintf("/titles/%d/translate", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
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
