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

func GetGetRecentlyUpdatedTitlesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":       GetRecentlyUpdatedTitlesSuccess(env),
		"invalid limit": GetRecentlyUpdatedTitlesWithInvalidLimit(env),
	}
}

func GetRecentlyUpdatedTitlesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, []string{"fighting", "action"})
			if err != nil {
				t.Fatal(err)
			}

			volumeID, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := testhelpers.CreateChapter(env.DB, volumeID, teamID, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/titles/recently-updated", h.GetRecentlyUpdatedTitles)

		req := httptest.NewRequest("GET", "/titles/recently-updated?limit=10", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) < 3 {
			t.Fatal("не все тайтлы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
		if _, ok := resp[0]["author"]; !ok {
			t.Fatal("автор не дошел")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не дошел")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("команда не дошла")
		}
		if genres, ok := resp[0]["genres"]; !ok || len(genres.([]any)) != 2 {
			t.Fatal("возникли проблемы с жанрами")
		}
	}
}

func GetRecentlyUpdatedTitlesWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/titles/recently-updated", h.GetRecentlyUpdatedTitles)

		req := httptest.NewRequest("GET", "/titles/recently-updated?limit=P_P", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
