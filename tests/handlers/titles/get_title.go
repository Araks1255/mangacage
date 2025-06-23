package titles

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":          GetTitleSuccess(env),
		"no views success": GetTitleWithNoViewsSuccess(env),
		"wrong id":         GetTitleWithWrongId(env),
		"invalid id":       GetTitleWithInvalidId(env),
	}
}

func GetTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, []string{"fighting", "action"}, []string{"maids", "japan"})
		if err != nil {
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

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/titles/:id", h.GetTitle)

		url := fmt.Sprintf("/titles/%d", titleID)
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
			t.Fatal("название не дошло")
		}
		if _, ok := resp["author"]; !ok {
			t.Fatal("автор не дошел")
		}
		if _, ok := resp["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
		if views, ok := resp["views"]; !ok || views.(float64) != 1 {
			t.Fatal("возникли проблемы с просмотрами", views)
		}
		if genres, ok := resp["genres"]; !ok || len(genres.([]any)) != 2 {
			t.Fatal("возникли проблемы с жанрами")
		}
		if tags, ok := resp["tags"]; !ok || len(tags.([]any)) != 2 {
			t.Fatal("возникли проблемы с тегами")
		}
	}
}

func GetTitleWithNoViewsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, []string{"fighting", "action"}, []string{"maids"})
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.GET("/titles/:id", h.GetTitle)

		url := fmt.Sprintf("/titles/%d", titleID)
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
			t.Fatal("название не дошло")
		}
		if _, ok := resp["englishName"]; !ok {
			t.Fatal("английское название не дошло")
		}
		if _, ok := resp["originalName"]; !ok {
			t.Fatal("оригинальное название не дошло")
		}
		if _, ok := resp["yearOfRelease"]; !ok {
			t.Fatal("год выпуска не дошел")
		}
		if _, ok := resp["ageLimit"]; !ok {
			t.Fatal("возрастное ограничение не дошло")
		}
		if _, ok := resp["type"]; !ok {
			t.Fatal("тип не дошел")
		}
		if _, ok := resp["translatingStatus"]; !ok {
			t.Fatal("статус перевода не дошёл")
		}
		if _, ok := resp["publishingStatus"]; !ok {
			t.Fatal("статус выпуска не дошёл")
		}
		if _, ok := resp["author"]; !ok {
			t.Fatal("автор не дошел")
		}
		if _, ok := resp["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
		if genres, ok := resp["genres"]; !ok || len(genres.([]any)) != 2 {
			t.Fatal("возникли проблемы с жанрами")
		}
		if tags, ok := resp["tags"]; !ok || len(tags.([]any)) != 1 || tags.([]any)[0].(string) != "maids" {
			t.Fatal("возникли проблемы с тегами")
		}
	}
}

func GetTitleWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil)

		titleID := 9223372036854775807

		r := gin.New()
		r.GET("/titles/:id", h.GetTitle)

		url := fmt.Sprintf("/titles/%d", titleID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitleWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil)

		invalidTitleID := "Y_Y"

		r := gin.New()
		r.GET("/titles/:id", h.GetTitle)

		url := fmt.Sprintf("/titles/%s", invalidTitleID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
