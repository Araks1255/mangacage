package tags

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/tags"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTagsScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":      GetTagsWithAllParamsSuccess(env),
		"success with query":      GetTagsSuccessWithQuery(env),
		"success with pagination": GetTagsWithPagination(env),
		"not found":               GetTagsNotFound(env),
	}
}

func GetTagsWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		_, err := testhelpers.CreateTags(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		h := tags.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/tags", h.GetTags)

		url := "/tags?page=1&limit=20"
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
			t.Fatal("неверное количество тегов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetTagsSuccessWithQuery(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		tagsIDs, err := testhelpers.CreateTags(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		var tagName string
		if err := env.DB.Raw("SELECT name FROM tags WHERE id = ?", tagsIDs[0]).Scan(&tagName).Error; err != nil {
			t.Fatal(err)
		}

		h := tags.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/tags", h.GetTags)

		url := fmt.Sprintf("/tags?query=%s", tagName)
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
			t.Fatal("неверное количество тегов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetTagsWithPagination(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		_, err := testhelpers.CreateTags(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		h := tags.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/tags", h.GetTags)

		tagsIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/tags?page=%d&limit=1&sort=createdAt", i)
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

			tagsIDs[i-1] = uint(id)
		}

		if tagsIDs[0]-tagsIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией")
		}
	}
}

func GetTagsNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := tags.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/tags", h.GetTags)

		req := httptest.NewRequest("GET", "/tags?query=nonexistenttag", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}
