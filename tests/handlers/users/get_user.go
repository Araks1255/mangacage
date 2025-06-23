package users

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetUserScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":    GetUserSuccess(env),
		"not found":  GetUserNotFound(env),
		"invalid id": GetUserInvalidId(env),
		"invisible":  GetUserInvisible(env),
	}
}

func GetUserSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{
			Visible: true,
			Roles:   []string{"typer", "translater"},
		})
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
		r.GET("/users/:id", h.GetUser)

		url := fmt.Sprintf("/users/%d", userID)
		req := httptest.NewRequest("GET", url, nil)

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

func GetUserNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users/:id", h.GetUser)

		req := httptest.NewRequest("GET", "/users/999999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetUserInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users/:id", h.GetUser)

		req := httptest.NewRequest("GET", "/users/T_T", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetUserInvisible(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{
			Visible: false,
			Roles:   []string{"typer"},
		})
		if err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users/:id", h.GetUser)

		url := fmt.Sprintf("/users/%d", userID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal("невидимый пользователь найден")
		}
	}
}
