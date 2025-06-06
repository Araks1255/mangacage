package favorites

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/favorites"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetFavoriteChaptersScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                   GetFavoriteChaptersSuccess(env),
		"unauthorized":              GetFavoriteChaptersUnauthorized(env),
		"without favorite chapters": GetFavoriteChaptersWithoutFavoriteChapters(env),
		"invalid limit":             GetFavoriteChaptersWithInvalidLimit(env),
	}
}

func GetFavoriteChaptersSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
			if err != nil {
				t.Fatal(err)
			}

			if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
				t.Fatal(err)
			}
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/chapters", h.GetFavoriteChapters)

		req := httptest.NewRequest("GET", "/users/me/favorites/chapters", nil)

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
			t.Fatalf(w.Body.String())
		}

		if len(resp) < 3 {
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

func GetFavoriteChaptersUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/chapters", h.GetFavoriteChapters)

		req := httptest.NewRequest("GET", "/users/me/favorites/chapters", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetFavoriteChaptersWithoutFavoriteChapters(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/chapters", h.GetFavoriteChapters)

		req := httptest.NewRequest("GET", "/users/me/favorites/chapters", nil)

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

func GetFavoriteChaptersWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/chapters", h.GetFavoriteChapters)

		req := httptest.NewRequest("GET", "/users/me/favorites/chapters?limit=ж_ж", nil)

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
