package search

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/search"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func SearchTeams(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		var teamName string
		env.DB.Raw("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)
		if len(teamName) < 5 {
			t.Fatal("не удалось получить название созданной команды")
		}

		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		query := teamName[:5]
		url := fmt.Sprintf("/search?type=teams&query=%s&limit=10", query)
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

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("не был получен id")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("не было получено название")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("не было получено время создания")
		}
		if _, ok := resp[0]["leader"]; !ok {
			t.Fatal("не был получен лидер")
		}
		if _, ok := resp[0]["leaderId"]; !ok {
			t.Fatal("не был получен id лидера")
		}
	}
}
