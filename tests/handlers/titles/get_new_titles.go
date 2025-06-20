package titles

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetNewTitlesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":       GetNewTitlesSuccess(env),
		"invalid limit": GetNewTitlesWithInvalidLimit(env),
	}
}

func GetNewTitlesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/titles/new", h.GetNewTitles)

		req := httptest.NewRequest("GET", "/titles/new?limit=5", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) < 2 {
			t.Fatal("не все тайтлы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetNewTitlesWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/titles/new", h.GetNewTitles)

		req := httptest.NewRequest("GET", "/titles/new?limit=*-*", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf(w.Body.String())
		}
	}
}
