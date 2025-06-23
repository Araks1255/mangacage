package teams

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTeamsScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":      GetTeamsWithAllParamsSuccess(env),
		"success with query":      GetTeamsSuccessWithQuery(env),
		"success with pagination": GetTeamsWithPagination(env),
		"not found":               GetTeamsNotFound(env),
		"invalid order":           GetTeamsWithInvalidOrder(env),
	}
}

func GetTeamsWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateTeam(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}
		}

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams", h.GetTeams)

		url := "/teams?sort=createdAt&order=desc&page=1&limit=20"
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

		if len(resp) < 2 {
			t.Fatal("неверное количество команд")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошёл")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetTeamsSuccessWithQuery(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var teamName string
		if err := env.DB.Raw("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName).Error; err != nil {
			t.Fatal(err)
		}

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams", h.GetTeams)

		url := fmt.Sprintf("/teams?query=%s", teamName)
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
			t.Fatal("неверное количество команд")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошёл")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetTeamsWithPagination(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateTeam(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}
		}

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams", h.GetTeams)

		teamsIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/teams?limit=1&page=%d&sort=createdAt", i)
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

			teamsIDs[i-1] = uint(id)
		}

		if teamsIDs[0]-teamsIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией")
		}
	}
}

func GetTeamsNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams", h.GetTeams)

		req := httptest.NewRequest("GET", "/teams?query=nonexistentteam", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTeamsWithInvalidOrder(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateTeam(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}
		}

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams", h.GetTeams)

		req := httptest.NewRequest("GET", "/teams?order=notvalid&sort=createdAt", nil)
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
			t.Fatal("возникли проблемы с количеством команд")
		}

		if uint(resp[0]["id"].(float64))-uint(resp[1]["id"].(float64)) != 1 { // При невалидном order должен выставиться desc
			t.Fatal("возникли проблемы с порядком команд")
		}
	}
}
