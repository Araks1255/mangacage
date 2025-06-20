package moderation

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyChapterOnModerationPageScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         GetMyChapterOnModerationPageSuccess(env),
		"unauthorized":    GetMyChapterOnModerationPageUnauthorized(env),
		"edited":          GetMyEditedChapterOnModerationPage(env),
		"other`s chapter": GetOthersChapterOnModerationPage(env),
		"wrong id":        GetMyChapterOnModerationPageWithWrongId(env),
		"wrong page":      GetMyChapterOnModerationPageWithWrongPage(env),
		"invalid id":      GetMyChapterOnModerationPageWithInvalidId(env),
		"invalid page":    GetMyChapterOnModerationPageWithInvalidPage(env),
	}
}

func GetMyChapterOnModerationPageSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		pages := make([][]byte, 1, 1)

		if pages[0], err = os.ReadFile("./test_data/test_chapter_page.png"); err != nil {
			t.Fatal(err)
		}

		chapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: pages, Collection: chaptersPages},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, chaptersPages, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d/page/0", chapterOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		if len(w.Body.Bytes()) != len(pages[0]) {
			t.Fatalf("возникли проблемы со страницей")
		}
	}
}

func GetMyChapterOnModerationPageUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters/18/page/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetMyEditedChapterOnModerationPage(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		editedChapterID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Edited: true},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, chaptersPages, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d/page/0", editedChapterID)
		req := httptest.NewRequest("GET", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetOthersChapterOnModerationPage(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		pages := make([][]byte, 1, 1)
		if pages[0], err = os.ReadFile("./test_data/test_chapter_page.png"); err != nil {
			t.Fatal(err)
		}

		otherUserChapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, otherUserID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: pages, Collection: chaptersPages},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, chaptersPages, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d/page/0", otherUserChapterOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyChapterOnModerationPageWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterOnModerationID := 9223372036854775807

		h := moderation.NewHandler(env.DB, nil, chaptersPages, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d/page/0", chapterOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetMyChapterOnModerationPageWithWrongPage(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, chaptersPages, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d/page/9223372036854775807", chapterOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetMyChapterOnModerationPageWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidChapterOnModerationID := "J_J"

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		url := fmt.Sprintf("/users/me/moderation/chapters/%s/page/0", invalidChapterOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetMyChapterOnModerationPageWithInvalidPage(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters/18/page/J_J", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
