package moderation

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyChaptersOnModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         GetMyChaptersOnModerationSuccess(env),
		"success new":     GetMyNewChaptersOnModerationSuccess(env),
		"success edited":  GetMyEditedChaptersOnModerationSuccess(env),
		"unauthorized":    GetMyChaptersOnModerationUnauthorized(env),
		"without volumes": GetMyChaptersOnModerationWithoutChapters(env),
		"invalid type":    GetMyChaptersOnModerationWithInvalidType(env),
		"invalid limit":   GetMyChaptersOnModerationWithInvalidLimit(env),
	}
}

func GetMyChaptersOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		blankPages := make([][]byte, 2, 2)

		if _, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: blankPages, Collection: chaptersPages},
		); err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
			env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Edited: true, Pages: blankPages, Collection: chaptersPages},
		); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

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
			t.Fatal("не все главы дошли")
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
		if _, ok := resp[0]["volume"]; !ok {
			t.Fatal("том не дошел")
		}
		if _, ok := resp[0]["volumeId"]; !ok {
			t.Fatal("id тома не дошел")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
	}
}

func GetMyNewChaptersOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		blankPages := make([][]byte, 2, 2)

		for i := 0; i < 2; i++ {
			if _, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
				env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: blankPages, Collection: chaptersPages},
			); err != nil {
				t.Fatal(err)
			}
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters?type=new", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) == 0 {
			t.Fatal("главы на модерации не дошли")
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
		if numberOfPages, ok := resp[0]["numberOfPages"]; !ok || numberOfPages.(float64) != 2 {
			t.Fatal("возникли проблемы с количеством страниц")
		}
		if _, ok := resp[0]["volume"]; !ok {
			t.Fatal("том не дошел")
		}
		if _, ok := resp[0]["volumeId"]; !ok {
			t.Fatal("id тома не дошел")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
	}
}

func GetMyEditedChaptersOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		blankPages := make([][]byte, 2, 2)

		for i := 0; i < 2; i++ {
			if _, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
				env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: blankPages, Collection: chaptersPages, Edited: true},
			); err != nil {
				t.Fatal(err)
			}
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters?type=edited", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatalf(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) == 0 {
			t.Fatal("главы на модерации не дошли")
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
		if _, ok := resp[0]["volume"]; !ok {
			t.Fatal("том не дошел")
		}
		if _, ok := resp[0]["volumeId"]; !ok {
			t.Fatal("id тома не дошел")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
	}
}

func GetMyChaptersOnModerationUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyChaptersOnModerationWithoutChapters(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters", nil)

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

func GetMyChaptersOnModerationWithInvalidType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters?type=O_O", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyChaptersOnModerationWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/chapters", h.GetMyChaptersOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/chapters?limit=o_O", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf(w.Body.String())
		}
	}
}
