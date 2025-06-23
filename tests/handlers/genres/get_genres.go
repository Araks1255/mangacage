package genres

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/genres"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetGenresScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":          GetGenresWithAllParamsSuccess(env),
		"success with query":          GetGenresWithQuerySuccess(env),
		"success with pagination":     GetGenresWithPaginationSuccess(env),
		"success my favorites":        GetMyFavoriteGenresSuccess(env),
		"success favorited by user":   GetGenresFavoritedByUserSuccess(env),
		"favorites unauthorized":      GetMyFavoriteGenresUnauthorized(env),
		"favorited by invisible user": GetGenresFavoritedByInvisibleUser(env),
		"not found":                   GetGenresNotFound(env),
	}
}

func GetGenresWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		_, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		url := "/genres?page=1&limit=20"
		req := httptest.NewRequest("GET", url, nil)

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
			t.Fatal("неверное количество жанров")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetGenresWithQuerySuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		genresIDs, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		var genreName string
		if err := env.DB.Raw("SELECT name FROM genres WHERE id = ?", genresIDs[0]).Scan(&genreName).Error; err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		url := fmt.Sprintf("/genres?query=%s", genreName)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 1 {
			t.Fatal("неверное количество жанров")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetGenresWithPaginationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		_, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		genresIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/genres?page=%d&limit=1&sort=createdAt", i)
			req := httptest.NewRequest("GET", url, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Fatal(w.Body.String())
			}

			var resp []map[string]any
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatal(err)
			}

			id, ok := resp[0]["id"].(float64)
			if !ok {
				t.Fatal("возникли проблемы с получением id")
			}

			genresIDs[i-1] = uint(id)
		}

		if genresIDs[0]-genresIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией")
		}
	}
}

func GetMyFavoriteGenresSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddGenreToFavorites(env.DB, userID, genresIDs[0]); err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		req := httptest.NewRequest("GET", "/genres?myFavorites=true", nil)

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

		if len(resp) != 1 {
			t.Fatal("неверное количество жанров")
		}
		if uint(resp[0]["id"].(float64)) != genresIDs[0] {
			t.Fatal("пришел не тот жанр")
		}
	}
}

func GetGenresFavoritedByUserSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Visible: true})
		if err != nil {
			t.Fatal(err)
		}

		genresIDs, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddGenreToFavorites(env.DB, userID, genresIDs[0]); err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		url := fmt.Sprintf("/genres?favoritedBy=%d", userID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 1 {
			t.Fatal("неверное количество жанров")
		}
		if uint(resp[0]["id"].(float64)) != genresIDs[0] {
			t.Fatal("пришел не тот жанр")
		}
	}
}

func GetMyFavoriteGenresUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		req := httptest.NewRequest("GET", "/genres?myFavorites=true", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetGenresFavoritedByInvisibleUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddGenreToFavorites(env.DB, userID, genresIDs[0]); err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		url := fmt.Sprintf("/genres?favoritedBy=%d", userID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetGenresNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := genres.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/genres", h.GetGenres)

		req := httptest.NewRequest("GET", "/genres?query=nonexistentgenre", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}
