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

func GetGetFavoriteTitlesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                 GetFavoriteTitlesSuccess(env),
		"unauthorized":            GetFavoriteTitlesUnauthorized(env),
		"without favorite titles": GetFavoriteTitlesWithoutFavoriteTitles(env),
		"invalid limit":           GetFavoriteTitlesWithInvalidLimit(env),
	}
}

func GetFavoriteTitlesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, []string{"action"}, nil)
			if err != nil {
				t.Fatal(err)
			}

			if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
				t.Fatal(err)
			}
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/titles", h.GetFavoriteTitles)

		req := httptest.NewRequest("GET", "/users/me/favorites/titles", nil)

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

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) < 3 {
			t.Fatal("не все тайтлы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошёл")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["author"]; !ok {
			t.Fatal("автор не дошёл")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошёл")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("команда не дошла")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не дошёл")
		}

		if genres, ok := resp[0]["genres"]; !ok || len(genres.([]any)) != 1 {
			t.Fatal("возникли проблемы с жанрами")
		}
	}
}

func GetFavoriteTitlesUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/titles", h.GetFavoriteTitles)

		req := httptest.NewRequest("GET", "/users/me/favorites/titles", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetFavoriteTitlesWithoutFavoriteTitles(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/titles", h.GetFavoriteTitles)

		req := httptest.NewRequest("GET", "/users/me/favorites/titles", nil)

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

func GetFavoriteTitlesWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/titles", h.GetFavoriteTitles)

		req := httptest.NewRequest("GET", "/users/me/favorites/titles?limit=O_O", nil)
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
