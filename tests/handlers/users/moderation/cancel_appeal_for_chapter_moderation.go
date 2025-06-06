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

func GetCancelAppealForChapterModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success new":    CancelAppealForNewChapterModerationSuccess(env),
		"success edited": CancelAppealForEditedChapterModerationSuccess(env),
		"other`s appeal": CancelOthersAppealForChapterModeration(env),
	}
}

func CancelAppealForNewChapterModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersOnModerationPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		pages := make([][]byte, 1)

		if pages[0], err = os.ReadFile("./test_data/chapter_page.png"); err != nil {
			t.Fatal(err)
		}

		newChapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: pages, Collection: chaptersPages},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, chaptersPages, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d", newChapterOnModerationID)
		req := httptest.NewRequest("DELETE", url, nil)

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
	}
}

func CancelAppealForEditedChapterModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		editedChapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Edited: true},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d", editedChapterOnModerationID)
		req := httptest.NewRequest("DELETE", url, nil)

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
	}
}

func CancelOthersAppealForChapterModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserChapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(env.DB, otherUserID)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/chapters/%d", otherUserChapterOnModerationID)
		req := httptest.NewRequest("DELETE", url, nil)

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
