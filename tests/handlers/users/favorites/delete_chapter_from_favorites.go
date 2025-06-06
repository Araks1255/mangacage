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

func GetDeleteChapterFromFavoritesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":          DeleteChapterFromFavoritesSuccess(env),
		"unauthorized":     DeleteChapterFromFavoritesUnauthorized(env),
		"wrong id":         DeleteChapterFromFavoritesWithWrongId(env),
		"not in favorites": DeleteChapterThatIsNotInFavoritesFromFavorites(env),
		"invalid ID":       DeleteChapterFromFavoritesWithInvalidId(env),
	}
}

func DeleteChapterFromFavoritesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/favorites/chapters/:id", h.DeleteChapterFromFavorites)

		url := fmt.Sprintf("/users/me/favorites/chapters/%d", chapterID)
		req := httptest.NewRequest("DELETE", url, nil)

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
	}
}

func DeleteChapterFromFavoritesUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/favorites/chapters/:id", h.DeleteChapterFromFavorites)

		req := httptest.NewRequest("DELETE", "/users/me/favorites/chapters/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func DeleteChapterFromFavoritesWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID := 9223372036854775807

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/favorites/chapters/:id", h.DeleteChapterFromFavorites)

		url := fmt.Sprintf("/users/me/favorites/chapters/%d", chapterID)
		req := httptest.NewRequest("DELETE", url, nil)

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

func DeleteChapterThatIsNotInFavoritesFromFavorites(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/favorites/chapters/:id", h.DeleteChapterFromFavorites)

		url := fmt.Sprintf("/users/me/favorites/chapters/%d", chapterID)
		req := httptest.NewRequest("DELETE", url, nil)

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

func DeleteChapterFromFavoritesWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := favorites.NewHandler(env.DB)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/favorites/chapters/:id", h.DeleteChapterFromFavorites)

		req := httptest.NewRequest("DELETE", "/users/me/favorites/chapters/J_J", nil)

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
