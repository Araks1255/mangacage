package users

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyProfileScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      GetMyProfileSuccess(env),
		"unauthorized": GetMyProfileUnauthorized(env),
	}
}

func GetMyProfileSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer", "translater"}})
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

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me", h.GetMyProfile)

		req := httptest.NewRequest("GET", "/users/me", nil)

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

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if _, ok := resp["id"]; !ok {
			t.Fatal("id не дошёл")
		}
		if _, ok := resp["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp["userName"]; !ok {
			t.Fatal("имя пользователя не дошло")
		}
		if _, ok := resp["team"]; !ok {
			t.Fatal("команда не дошла")
		}
		if _, ok := resp["teamId"]; !ok {
			t.Fatal("id команды не дошёл")
		}
		if roles, ok := resp["roles"]; !ok || len(roles.([]any)) != 2 {
			t.Fatal("возникли проблемы с ролями")
		}
	}
}

func GetMyProfileUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me", h.GetMyProfile)

		req := httptest.NewRequest("GET", "/users/me", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}
