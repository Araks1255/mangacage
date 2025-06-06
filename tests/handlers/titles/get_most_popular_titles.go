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

func GetGetMostPopularTitlesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":       GetMostPopularTitlesSuccess(env),
		"invalid limit": GetMostPopularTitlesWithInvalidLimit(env),
	}
}

func GetMostPopularTitlesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 4; i++ {
			titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID, "fighting")
			if err != nil {
				t.Fatal(err)
			}

			if err := testhelpers.TranslateTitle(env.DB, teamID, titleID); err != nil {
				t.Fatal(err)
			}

			volumeID, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID)
			if err != nil {
				t.Fatal(err)
			}

			chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, teamID, userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := testhelpers.ViewChapter(env.DB, userID, chapterID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/titles/most-popular", h.GetMostPopularTitles)

		req := httptest.NewRequest("GET", "/titles/most-popular?limit=10", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) < 4 {
			t.Fatal("не все тайтлы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название тайтла не дошло")
		}
		if _, ok := resp[0]["author"]; !ok {
			t.Fatal("автор не дошел")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("команда не дошла")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не дошел")
		}
		if genres, ok := resp[0]["genres"]; !ok || len(genres.([]any)) == 0 || genres.([]any)[0].(string) != "fighting" {
			t.Fatal("возникли проблемы с жанрами")
		}
		if views, ok := resp[0]["views"]; !ok || views.(float64) != 1 {
			t.Fatal("возникли проблемы с просмотрами")
		}
	}
}

func GetMostPopularTitlesWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/titles/most-popular", h.GetMostPopularTitles)

		req := httptest.NewRequest("GET", "/titles/most-popular?limit=U_U", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
