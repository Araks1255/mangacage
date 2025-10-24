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

func GetCancelTitleTranslateRequestScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":              CancelTitleTranslateRequestSuccess(env),
		"another team request": CancelAnotherTeamTitleTranslateRequest(env),
	}
}

func CancelTitleTranslateRequestSuccess(env testenv.Env) func(*testing.T) {
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

		requestID, err := testhelpers.CreateTitleTranslateRequest(env.DB, titleID, teamID)
		if err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/titles/translate-requests/:id", h.CancelTitleTranslateRequest)

		url := fmt.Sprintf("/titles/translate-requests/%d", requestID)
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

func CancelAnotherTeamTitleTranslateRequest(env testenv.Env) func(*testing.T) {
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

		anotherTeamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		requestID, err := testhelpers.CreateTitleTranslateRequest(env.DB, titleID, anotherTeamID)
		if err != nil {
			t.Fatal(err)
		}

		h := translaterequests.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/titles/translate-requests/:id", h.CancelTitleTranslateRequest)

		url := fmt.Sprintf("/titles/translate-requests/%d", requestID)
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
