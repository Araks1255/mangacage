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

func SearchAuthors(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		authorID, err := testhelpers.CreateAuthor(env.DB, testhelpers.CreateAuthorOptions{Genres: []string{"fighting"}})
		if err != nil {
			t.Fatal(err)
		}

		var authorName string
		env.DB.Raw("SELECT name FROM authors WHERE id = ?", authorID).Scan(&authorName)
		if len(authorName) < 5 {
			t.Fatal("не удалось получить имя созданного автора")
		}

		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		query := authorName[:5]
		url := fmt.Sprintf("/search?type=authors&query=%s&limit=10", query)
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
			t.Fatal("не был отправлен id")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("не было отправлено имя")
		}
		if genres, ok := resp[0]["genres"]; !ok || len(genres.([]any)) != 1 {
			t.Fatal("не были отправлены жанры")
		}
	}
}
