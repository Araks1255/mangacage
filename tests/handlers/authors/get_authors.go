package authors

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/authors"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetAuthorsScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":      GetAuthorsWithAllParamsSuccess(env),
		"success with query":      GetAuthorsSuccessWithQuery(env),
		"success with pagination": GetAuthorsWithPagination(env),
		"not found":               GetAuthorsNotFound(env),
	}
}

func GetAuthorsWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateAuthor(env.DB); err != nil {
				t.Fatal(err)
			}
		}

		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors", h.GetAuthors)

		url := "/authors?page=1&limit=20"
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
			t.Fatal("неверное количество авторов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("имя не дошло")
		}
		if _, ok := resp[0]["englishName"]; !ok {
			t.Fatal("имя на английском не дошло")
		}
		if _, ok := resp[0]["originalName"]; !ok {
			t.Fatal("оригинальное имя не дошло")
		}
	}
}

func GetAuthorsSuccessWithQuery(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateAuthor(env.DB); err != nil {
				t.Fatal(err)
			}
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var authorName string
		if err := env.DB.Raw("SELECT name FROM authors WHERE id = ?", authorID).Scan(&authorName).Error; err != nil {
			t.Fatal(err)
		}

		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors", h.GetAuthors)

		url := fmt.Sprintf("/authors?query=%s", authorName)
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
			t.Fatal("неверное количество авторов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("имя не дошло")
		}
		if _, ok := resp[0]["englishName"]; !ok {
			t.Fatal("имя на английском не дошло")
		}
		if _, ok := resp[0]["originalName"]; !ok {
			t.Fatal("оригинальное имя не дошло")
		}
	}
}

func GetAuthorsWithPagination(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateAuthor(env.DB); err != nil {
				t.Fatal(err)
			}
		}

		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors", h.GetAuthors)

		authorsIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/authors?page=%d&limit=1&sort=createdAt", i)
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

			authorsIDs[i-1] = uint(id)
		}

		if authorsIDs[0]-authorsIDs[1] != 1 {
			var authors []map[string]any
			env.DB.Raw("SELECT * FROM authors").Scan(&authors)
			t.Fatal("возникли проблемы с пагинацией", authorsIDs, "\n", authors)
		}
	}
}

func GetAuthorsNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors", h.GetAuthors)

		req := httptest.NewRequest("GET", "/authors?query=nonexistentauthor", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}
