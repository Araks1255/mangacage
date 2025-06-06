package search

import (
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/search"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetSearchScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"authors":              SearchAuthors(env),
		"chapters":             SearchChapters(env),
		"teams":                SearchTeams(env),
		"titles":               SearchTitles(env),
		"volumes":              SearchVolumes(env),
		"wrong searching type": SearchWithWrongSearchingType(env),
		"no searching type":    SearchWithNoSearchingType(env),
		"no query":             SearchWithNoQuery(env),
		"invalid limit":        SearchWithInvalidLimit(env),
	}
}

func SearchWithWrongSearchingType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		req := httptest.NewRequest("GET", "/search?type=T_T&query=Araki&limit=10", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func SearchWithNoSearchingType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		req := httptest.NewRequest("GET", "/search?query=Araki&limit=10", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func SearchWithNoQuery(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		req := httptest.NewRequest("GET", "/search?type=authors&limit=10", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func SearchWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := search.NewHandler(env.DB)

		r := gin.New()
		r.GET("/search", h.Search)

		req := httptest.NewRequest("GET", "/search?type=authors&query=Araki&limit=O_O", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
