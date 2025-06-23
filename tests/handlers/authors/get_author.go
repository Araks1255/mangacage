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

func GetGetAuthorScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":    GetAuthorSuccess(env),
		"not found":  GetAuthorNotFound(env),
		"invalid id": GetAuthorInvalidId(env),
	}
}

func GetAuthorSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors/:id", h.GetAuthor)

		url := fmt.Sprintf("/authors/%d", authorID)
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
			t.Fatal("id не дошел")
		}
		if _, ok := resp["name"]; !ok {
			t.Fatal("имя не дошло")
		}
		if _, ok := resp["englishName"]; !ok {
			t.Fatal("имя на английском не дошло")
		}
		if _, ok := resp["originalName"]; !ok {
			t.Fatal("оригинальное имя не дошло")
		}
	}
}

func GetAuthorNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors/:id", h.GetAuthor)

		req := httptest.NewRequest("GET", "/authors/999999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetAuthorInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := authors.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/authors/:id", h.GetAuthor)

		req := httptest.NewRequest("GET", "/authors/invalid_id", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
