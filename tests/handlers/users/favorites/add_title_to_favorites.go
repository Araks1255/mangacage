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

func GetAddTitleToFavoritesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      AddTitleToFavoritesSuccess(env),
		"unauthorized": AddTitleToFavoritesUnauthorized(env),
		"twice":        AddTitleToFavoritesTwice(env),
		"wrong id":     AddTitleToFavoritesWithWrongId(env),
		"invalid id":   AddTitleToFavoritesWithInvalidId(env),
	}
}

func AddTitleToFavoritesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID, "fighting")
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/titles/:id", h.AddTitleToFavorites)

		url := fmt.Sprintf("/users/me/favorites/titles/%d", titleID)
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

func AddTitleToFavoritesUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/titles/:id", h.AddTitleToFavorites)

		req := httptest.NewRequest("POST", "/users/me/favorites/titles/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func AddTitleToFavoritesTwice(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID, "fighting")
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/titles/:id", h.AddTitleToFavorites)

		url := fmt.Sprintf("/users/me/favorites/titles/%d", titleID)
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

func AddTitleToFavoritesWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID := 9223372036854775807

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/titles/:id", h.AddTitleToFavorites)

		url := fmt.Sprintf("/users/me/favorites/titles/%d", titleID)
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

func AddTitleToFavoritesWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidTitleID := "()_()"

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/favorites/titles/:id", h.AddTitleToFavorites)

		url := fmt.Sprintf("/users/me/favorites/titles/%s", invalidTitleID)
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
