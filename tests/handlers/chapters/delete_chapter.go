package chapters

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetDeleteChapterScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":                  DeleteChapterSuccess(env),
		"unauthorized":             DeleteChapterByUnauthorizedUser(env),
		"non team leader":          DeleteChapterByNonTeamLeader(env),
		"does not translate title": DeleteChapterByUserWhoseTeamDoesNotTranslateTitle(env),
		"wrong id":                 DeleteChapterWithWrongId(env),
		"invalid id":               DeleteChapterWithInvalidId(env),
	}
}

func DeleteChapterSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, nil, chaptersPages)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/chapters/:id", h.DeleteChapter)

		url := fmt.Sprintf("/chapters/%d", chapterID)
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

func DeleteChapterByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/chapters/:id", h.DeleteChapter)

		req := httptest.NewRequest("DELETE", "/chapters/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteChapterByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/chapters/:id", h.DeleteChapter)

		url := fmt.Sprintf("/chapters/%d", chapterID)
		req := httptest.NewRequest("DELETE", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteChapterByUserWhoseTeamDoesNotTranslateTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/chapters/:id", h.DeleteChapter)

		url := fmt.Sprintf("/chapters/%d", chapterID)
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

func DeleteChapterWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID := 9223372036854775807

		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/chapters/:id", h.DeleteChapter)

		url := fmt.Sprintf("/chapters/%d", chapterID)
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

func DeleteChapterWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID := "._."

		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/chapters/:id", h.DeleteChapter)

		url := fmt.Sprintf("/chapters/%s", chapterID)
		req := httptest.NewRequest("DELETE", url, nil)

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
