package moderation

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyTeamOnModerationScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":        GetMyTeamOnModerationSuccess(env),
		"success new":    GetMyNewTeamOnModerationSuccess(env),
		"success edited": GetMyEditedTeamOnModerationSuccess(env),
		"unauthorized":   GetMyTeamOnModerationUnauthorized(env),
		"without team":   GetMyTeamOnModerationWithoutTeam(env),
	}
}

func GetMyTeamOnModerationSuccess(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateTeamOnModeration(env.DB, userID); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/team", h.GetMyTeamOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/team", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf(w.Body.String())
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if _, ok := resp["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetMyNewTeamOnModerationSuccess(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateTeamOnModeration(env.DB, userID); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/team", h.GetMyTeamOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/team?type=new", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf(w.Body.String())
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if _, ok := resp["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetMyEditedTeamOnModerationSuccess(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		existingTeamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, userID, existingTeamID); err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateTeamOnModeration(
			env.DB, userID, moderationHelpers.CreateTeamOnModerationOptions{ExistingID: existingTeamID},
		); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/team", h.GetMyTeamOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/team?type=edited", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf(w.Body.String())
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if _, ok := resp["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp["existing"]; !ok {
			t.Fatal("оригинальная команда не дошла")
		}
		if _, ok := resp["existingId"]; !ok {
			t.Fatal("id оригинальной команды не дошел")
		}
	}
}

func GetMyTeamOnModerationUnauthorized(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/team", h.GetMyTeamOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/team", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyTeamOnModerationWithoutTeam(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/team", h.GetMyTeamOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/team", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatalf(w.Body.String())
		}
	}
}
