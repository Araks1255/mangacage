package chapters

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetChapterScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":    GetChapterSuccess(env),
		"wrong id":   GetChapterWithWrongId(env),
		"invalid id": GetChapterWithInvalidId(env),
	}
}

func GetChapterSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/chapters/:id", h.GetChapter)

		url := fmt.Sprintf("/chapters/%d", chapterID)
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
		if _, ok := resp["volume"]; !ok {
			t.Fatal("в ответе нет тома")
		}
		if _, ok := resp["title"]; !ok {
			t.Fatal("в ответе нет названия тайтла")
		}
		if _, ok := resp["titleId"]; !ok {
			t.Fatal("в ответе нет id тайтла")
		}
		if _, ok := resp["team"]; !ok {
			t.Fatal("в ответе нет команды")
		}
		if _, ok := resp["teamId"]; !ok {
			t.Fatal("в ответе нет id команды")
		}
	}
}

func GetChapterWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chapterID := 9223372036854775807

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/chapters/:id", h.GetChapter)

		url := fmt.Sprintf("/chapters/%d", chapterID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChapterWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chapterID := "^-^"

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/chapters/:id", h.GetChapter)

		url := fmt.Sprintf("/chapters/%s", chapterID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
