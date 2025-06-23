package genres

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/handlers/genres"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAddGenreScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                     AddGenreSuccess(env),
		"with the same name as genre": AddGenreWithTheSameNameAsGenre(env),
		"with the same name as genre on moderation": AddGenreWithTheSameNameAsGenreOnModeration(env),
	}
}

func AddGenreSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		body := models.GenreOnModerationDTO{
			Name: uuid.New().String(),
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/genres", h.AddGenre)

		req := httptest.NewRequest("POST", "/genres", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func AddGenreWithTheSameNameAsGenre(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs, err := testhelpers.CreateGenres(env.DB, 1)
		if err != nil {
			t.Fatal(err)
		}

		var genreName string
		if err := env.DB.Raw("SELECT name FROM genres WHERE id = ?", genresIDs[0]).Scan(&genreName).Error; err != nil {
			t.Fatal(err)
		}

		body := models.GenreOnModerationDTO{
			Name: genreName,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/genres", h.AddGenre)

		req := httptest.NewRequest("POST", "/genres", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func AddGenreWithTheSameNameAsGenreOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genreOnModerationID, err := moderation.CreateGenreOnModeration(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var genreOnModerationName string
		if err := env.DB.Raw("SELECT name FROM genres WHERE id = ?", genreOnModerationID).Scan(&genreOnModerationName).Error; err != nil {
			t.Fatal(err)
		}

		body := models.GenreOnModerationDTO{
			Name: genreOnModerationName,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := genres.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/genres", h.AddGenre)

		req := httptest.NewRequest("POST", "/genres", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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
