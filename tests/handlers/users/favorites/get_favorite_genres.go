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

func GetGetFavoriteGenresScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                 GetFavoriteGenresSuccess(env),
		"unauthorized":            GetFavoriteGenresUnauthorized(env),
		"without favorite genres": GetFavoriteGenresWithoutFavoriteGenres(env),
		"invalid limit":           GetFavoriteGenresWithInvalidLimit(env),
	}
}

func GetFavoriteGenresSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []uint
		if err := env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs).Error; err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := testhelpers.AddGenreToFavorites(env.DB, userID, genresIDs[i]); err != nil {
				t.Fatal(err)
			}
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/genres", h.GetFavoriteGenres)

		req := httptest.NewRequest("GET", "/users/me/favorites/genres", nil)

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

		if len(resp) < 2 {
			t.Fatal("не все жанры дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetFavoriteGenresUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/genres", h.GetFavoriteGenres)

		req := httptest.NewRequest("GET", "/users/me/favorites/genres", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetFavoriteGenresWithoutFavoriteGenres(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/genres", h.GetFavoriteGenres)

		req := httptest.NewRequest("GET", "/users/me/favorites/genres", nil)

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

func GetFavoriteGenresWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/favorites/genres", h.GetFavoriteGenres)

		req := httptest.NewRequest("GET", "/users/me/favorites/genres?limit=-_-", nil)

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
