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

func SearchChapters(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var chapterName string
		env.DB.Raw("SELECT name FROM chapters WHERE id = ?", chapterID).Scan(&chapterName)
		if len(chapterName) < 5 {
			t.Fatal("не удалось найти название созданной главы")
		}

		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		query := chapterName[:5]
		url := fmt.Sprintf("/search?type=chapters&query=%s&limit=10", query)
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
		if _, ok := resp[0]["volume"]; !ok {
			t.Fatal("не был получен том")
		}
		if _, ok := resp[0]["volumeId"]; !ok {
			t.Fatal("не был получен id тома")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("не был получен тайтл")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("не был получен id тайтла")
		}
	}
}
