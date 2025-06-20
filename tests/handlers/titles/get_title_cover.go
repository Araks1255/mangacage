package titles

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTitleCoverScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":    GetTitleCoverSuccess(env),
		"invalid id": GetTitleCoverWithInvalidId(env),
		"wrong id":   GetTitleCoverWithWrongId(env),
	}
}

func GetTitleCoverSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		cover, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Cover: cover, Collection: titlesCovers})
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, nil, titlesCovers)

		r := gin.New()
		r.GET("/titles/:id/cover", h.GetTitleCover)

		url := fmt.Sprintf("/titles/%d/cover", titleID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		if len(w.Body.Bytes()) != len(cover) {
			t.Fatal("возникли проблемы с обложкой")
		}
	}
}

func GetTitleCoverWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		h := titles.NewHandler(env.DB, nil, titlesCovers)

		titleID := 9223372036854775807

		r := gin.New()
		r.GET("/titles/:id/cover", h.GetTitleCover)

		url := fmt.Sprintf("/titles/%d/cover", titleID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitleCoverWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil)

		invalidTitleID := "::_::"

		r := gin.New()
		r.GET("/titles/:id/cover", h.GetTitleCover)

		url := fmt.Sprintf("/titles/%s/cover", invalidTitleID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
