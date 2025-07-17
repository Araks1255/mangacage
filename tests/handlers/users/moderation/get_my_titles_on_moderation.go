package moderation

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyTitlesOnModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":        GetMyTitlesOnModerationSuccess(env),
		"success new":    GetMyNewTitlesOnModerationSuccess(env),
		"success edited": GetMyEditedTitlesOnModerationSuccess(env),
		"unauthorized":   GetMyTitlesOnModerationUnauthorized(env),
		"without titles": GetMyTitlesOnModerationWithoutTitles(env),
		"invalid type":   GetMyTitlesOnModerationWithInvalidType(env),
		"invalid limit":  GetMyTitlesOnModerationWithInvalidLimit(env),
	}
}

func GetMyTitlesOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateTitleOnModeration(
			env.DB, userID,
			moderationHelpers.CreateTitleOnModerationOptions{
				Genres:   []string{"action", "fighting"},
				Tags:     []string{"maids", "japan"},
				AuthorID: authorID,
			},
		); err != nil {
			t.Fatal(err)
		}

		existingTitleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateTitleOnModeration(
			env.DB, userID, moderationHelpers.CreateTitleOnModerationOptions{
				ExistingID: existingTitleID, AuthorID: authorID,
				Genres: []string{"action", "fighting"},
				Tags:   []string{"maids", "japan"},
			},
		); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles", nil)

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
		if _, ok := resp[0]["author"]; !ok {
			t.Fatal("автор не дошел")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
	}
}

func GetMyNewTitlesOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(userID, "<===")

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := moderationHelpers.CreateTitleOnModeration(
				env.DB, userID,
				moderationHelpers.CreateTitleOnModerationOptions{
					Genres:   []string{"action", "fighting"},
					Tags:     []string{"maids", "japan"},
					AuthorID: authorID,
				},
			); err != nil {
				t.Fatal(err)
			}
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles?moderationType=new", nil)

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
		if _, ok := resp[0]["author"]; !ok {
			t.Fatal("автор не дошел")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
	}
}

func GetMyEditedTitlesOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			existingTitleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := moderationHelpers.CreateTitleOnModeration(
				env.DB, userID, moderationHelpers.CreateTitleOnModerationOptions{
					ExistingID: existingTitleID,
					Genres:     []string{"action", "fighting"},
					AuthorID:   authorID,
				},
			); err != nil {
				t.Fatal(err)
			}
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles?moderationType=edited", nil)

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
		if _, ok := resp[0]["existing"]; !ok {
			t.Fatal("оригинальный тайтл не дошел")
		}
		if _, ok := resp[0]["existingId"]; !ok {
			t.Fatal("id оригинального тайтла не дошел")
		}
	}
}

func GetMyTitlesOnModerationUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyTitlesOnModerationWithoutTitles(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles", nil)

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

func GetMyTitlesOnModerationWithInvalidType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles?moderationType=р_q", nil)

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

func GetMyTitlesOnModerationWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles", h.GetMyTitlesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles?limit=T_T", nil)

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
