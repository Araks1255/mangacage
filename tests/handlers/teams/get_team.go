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

func GetGetTeamScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":    GetTeamSuccess(env),
		"wrong id":   GetTeamWithWrongId(env),
		"invalid id": GetTeamWithInvalidId(env),
	}
}

func GetTeamSuccess(env testenv.Env) func(*testing.T) {
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

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams/:id", h.GetTeam)

		url := fmt.Sprintf("/teams/%d", teamID)
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
			t.Fatal("в ответе нет id")
		}
		if _, ok := resp["name"]; !ok {
			t.Fatal("в ответе нет названия")
		}
		if _, ok := resp["createdAt"]; !ok {
			t.Fatal("в ответе нет времени создания")
		}
		if _, ok := resp["leader"]; !ok {
			t.Fatal("в ответе нет имени лидера")
		}
		if _, ok := resp["leaderId"]; !ok {
			t.Fatal("в ответе нет id лидера")
		}
	}
}

func GetTeamWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamID := 9223372036854775807

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams/:id", h.GetTeam)

		url := fmt.Sprintf("/teams/%d", teamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTeamWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		invalidTeamID := "!_!"

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/teams/:id", h.GetTeam)

		url := fmt.Sprintf("/teams/%s", invalidTeamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
