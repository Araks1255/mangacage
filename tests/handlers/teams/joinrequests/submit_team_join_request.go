package joinrequests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetSubmitTeamJoinRequestScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":              SubmitTeamJoinRequestSuccess(env),
		"unauthorized":         SubmitTeamJoinRequestByUnauthorizedUser(env),
		"invalid team id":      SubmitTeamJoinRequestWithInvalidTeamId(env),
		"user already in team": SubmitTeamJoinRequestByUserThatAlreadyInTeam(env),
		"wrong role id":        SubmitTeamJoinRequestWithWrongRoleId(env),
		"repeated":             SubmitTeamJoinRequestTwice(env),
		"wrong team id":        SubmitTeamJoinRequestWithWrongTeamId(env),
	}
}

func SubmitTeamJoinRequestSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var roleID uint
		env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID)
		if roleID == 0 {
			t.Fatal("роль не найдена")
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		body := gin.H{
			"introductoryMessage": "message",
			"roleId":              roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/%d/join-requests", teamID)
		req := httptest.NewRequest("POST", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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
	}
}

func SubmitTeamJoinRequestByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		req := httptest.NewRequest("POST", "/teams/18/join-requests", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func SubmitTeamJoinRequestWithInvalidTeamId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidTeamID := ";|"

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		url := fmt.Sprintf("/teams/%s/join-requests", invalidTeamID)
		req := httptest.NewRequest("POST", url, nil)

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

func SubmitTeamJoinRequestByUserThatAlreadyInTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		anotherTeamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, userID, anotherTeamID); err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		url := fmt.Sprintf("/teams/%d/join-requests", teamID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")

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

func SubmitTeamJoinRequestWithWrongRoleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		roleID := 9223372036854775807

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		body := gin.H{
			"introductoryMessage": "message",
			"roleId":              roleID,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/teams/%d/join-requests", teamID)
		req := httptest.NewRequest("POST", url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func SubmitTeamJoinRequestTwice(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := testhelpers.CreateTeamJoinRequest(env.DB, userID, teamID, ""); err != nil {
			t.Fatal(err)
		}

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		url := fmt.Sprintf("/teams/%d/join-requests", teamID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")

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

func SubmitTeamJoinRequestWithWrongTeamId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID := 9223372036854775807

		h := joinrequests.NewHandler(env.DB, env.SecretKey, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

		url := fmt.Sprintf("/teams/%d/join-requests", teamID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")

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
