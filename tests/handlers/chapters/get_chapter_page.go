package chapters

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetChapterPageScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":                GetChapterPageSuccess(env),
		"wrong chapter id":       GetChapterPageWithWrongId(env),
		"invalid chapter id":     GetChapterPageWithInvalidId(env),
		"invalid number of page": GetChapterPageWithInvalidPage(env),
	}
}

func GetChapterPageSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		pages := make([][]byte, 1)
		pages[0], err = os.ReadFile("./test_data/chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID, testhelpers.CreateChapterOptions{Pages: pages, Collection: chaptersPages})
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, chaptersPages)

		r := gin.New()
		r.GET("/chapters/:id/page/:page", h.GetChapterPage)

		url := fmt.Sprintf("/chapters/%d/page/0", chapterID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		if len(w.Body.Bytes()) != len(pages[0]) {
			t.Fatal("фото не отправилось")
		}
	}
}

func GetChapterPageWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		wrongChapterID := 9223372036854775807

		h := chapters.NewHandler(env.DB, nil, chaptersPages)

		r := gin.New()
		r.GET("/chapters/:id/page/:page", h.GetChapterPage)

		url := fmt.Sprintf("/chapters/%d/page/0", wrongChapterID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChapterPageWithInvalidId(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil)

		invalidChapterID := "o_O"

		r := gin.New()
		r.GET("/chapters/:id/page/:page", h.GetChapterPage)

		url := fmt.Sprintf("/chapters/%s/page/0", invalidChapterID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChapterPageWithWrongPage(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil)

		page := 9223372036854775807

		r := gin.New()
		r.GET("/chapters/:id/page/:page", h.GetChapterPage)

		url := fmt.Sprintf("/chapters/0/page/%d", page)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChapterPageWithInvalidPage(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil)

		invalidPage := ":-)"

		r := gin.New()
		r.GET("/chapters/:id/page/:page", h.GetChapterPage)

		url := fmt.Sprintf("/chapters/0/page/%s", invalidPage)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
