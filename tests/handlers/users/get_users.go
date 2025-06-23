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

func GetGetUsersScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":      GetUsersWithAllParamsSuccess(env),
		"success with query":      GetUsersSuccessWithQuery(env),
		"success with pagination": GetUsersWithPagination(env),
		"not found":               GetUsersNotFound(env),
		"invalid order":           GetUsersWithInvalidOrder(env),
		"non visible":             GetUsersNonVisible(env),
	}
}

func GetUsersWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{
				TeamID:  teamID,
				Visible: true,
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users", h.GetUsers)

		url := fmt.Sprintf(
			"/users?sort=createdAt&order=desc&page=1&limit=20&teamId=%d",
			teamID,
		)

		req := httptest.NewRequest("GET", url, nil)

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
			t.Fatal("неверное количество пользователей")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошёл")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["userName"]; !ok {
			t.Fatal("имя пользователя не дошло")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не дошёл")
		}
	}
}

func GetUsersSuccessWithQuery(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Visible: true})
		if err != nil {
			t.Fatal(err)
		}

		var userName string
		if err := env.DB.Raw("SELECT user_name FROM users WHERE id = ?", userID).Scan(&userName).Error; err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users", h.GetUsers)

		url := fmt.Sprintf("/users?query=%s", userName)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 1 {
			t.Fatal("неверное количество пользователей")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошёл")
		}
		if _, ok := resp[0]["userName"]; !ok {
			t.Fatal("имя пользователя не дошло")
		}
	}
}

func GetUsersWithPagination(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Visible: true})
			if err != nil {
				t.Fatal(err)
			}
		}

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users", h.GetUsers)

		usersIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/users?limit=1&page=%d&sort=createdAt", i)
			req := httptest.NewRequest("GET", url, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Fatal(w.Body.String())
			}

			var resp []map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatal(err)
			}

			id, ok := resp[0]["id"].(float64)
			if !ok {
				t.Fatal("возникли проблемы с получением id")
			}

			usersIDs[i-1] = uint(id)
		}

		if usersIDs[0]-usersIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией")
		}
	}
}

func GetUsersNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users", h.GetUsers)

		req := httptest.NewRequest("GET", "/users?query=nonexistentuser", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetUsersWithInvalidOrder(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Visible: true})
			if err != nil {
				t.Fatal(err)
			}
		}

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users", h.GetUsers)

		req := httptest.NewRequest("GET", "/users?order=-_-&sort=createdAt", nil)
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
			t.Fatal("возникли проблемы с количеством пользователей")
		}

		if uint(resp[0]["id"].(float64))-uint(resp[1]["id"].(float64)) != 1 { // При невалидном order должен выставиться desc
			t.Fatal("возникли проблемы с порядком пользователей")
		}
	}
}

func GetUsersNonVisible(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var userName string
		if err := env.DB.Raw("SELECT user_name FROM users WHERE id = ?", userID).Scan(&userName).Error; err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/users", h.GetUsers)

		url := fmt.Sprintf("/users?query=%s", userName)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal("невидимый юзер нащелся")
		}
	}
}
