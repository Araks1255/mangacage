package favorites

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/favorites"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetAddGenreToFavoritesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      AddGenreToFavoritesSuccess(env),
		"unauthorized": AddGenreToFavoritesUnauthorized(env),
		"twice":        AddGenreToFavoritesTwice(env),
		"wrong id":     AddGenreToFavoritesWithWrongId(env),
		"invalid id":   AddGenreToFavoritesWithInvalidId(env),
	}
}

func AddGenreToFavoritesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genreID uint
		if err := env.DB.Raw("SELECT id FROM genres WHERE name = 'fighting'").Scan(&genreID).Error; err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/genres/:id", h.AddGenreToFavorites)

		url := fmt.Sprintf("/users/me/favorites/genres/%d", genreID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddGenreToFavoritesUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/genres/:id", h.AddGenreToFavorites)

		req := httptest.NewRequest("POST", "/users/me/favorites/genres/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddGenreToFavoritesTwice(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genreID uint
		if err := env.DB.Raw("SELECT id FROM genres WHERE name = 'fighting'").Scan(&genreID).Error; err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddGenreToFavorites(env.DB, userID, genreID); err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/genres/:id", h.AddGenreToFavorites)

		url := fmt.Sprintf("/users/me/favorites/genres/%d", genreID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 409 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddGenreToFavoritesWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genreID := 9223372036854775807

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/genres/:id", h.AddGenreToFavorites)

		url := fmt.Sprintf("/users/me/favorites/genres/%d", genreID)
		req := httptest.NewRequest("POST", url, nil)

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

func AddGenreToFavoritesWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidGenreID := "^-^"

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/genres/:id", h.AddGenreToFavorites)

		url := fmt.Sprintf("/users/me/favorites/genres/%s", invalidGenreID)
		req := httptest.NewRequest("POST", url, nil)

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
