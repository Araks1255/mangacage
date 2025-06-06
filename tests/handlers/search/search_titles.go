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

func SearchTitles(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, []string{"fighting"})
		if err != nil {
			t.Fatal(err)
		}

		var titleName string
		env.DB.Raw("SELECT name FROM titles WHERE id = ?", titleID).Scan(&titleName)
		if len(titleName) < 5 {
			t.Fatal("не удалось получить название созданного тайтла")
		}

		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		query := titleName[:5]
		url := fmt.Sprintf("/search?type=titles&query=%s&limit=10", query)
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
		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("не был получен id тайтла")
		}
		if _, ok := resp[0]["author"]; !ok {
			t.Fatal("не был получен автор")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("не был получен id автора")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("не была получена команда")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("не был получен id команды")
		}
		if genres, ok := resp[0]["genres"]; !ok || len(genres.([]any)) == 0 {
			t.Fatal("не были получены жанры")
		}
	}
}
