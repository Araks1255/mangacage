package authors

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/handlers/authors"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAddAuthorScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                      AddAuthorSuccess(env),
		"with the same name as author": AddAuthorWithTheSameNameAsAuthor(env),
		"with the same name as author on moderation": AddAuthorWithTheSameNameAsAuthorOnModeration(env),
	}
}

func AddAuthorSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		body := dto.CreateAuthorDTO{
			Name:         uuid.New().String(),
			EnglishName:  uuid.New().String(),
			OriginalName: "テストオーサー",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := authors.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/authors", h.AddAuthor)

		req := httptest.NewRequest("POST", "/authors", bytes.NewBuffer(jsonBody))
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

func AddAuthorWithTheSameNameAsAuthor(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var authorName string
		if err := env.DB.Raw("SELECT name FROM authors WHERE id = ?", authorID).Scan(&authorName).Error; err != nil {
			t.Fatal(err)
		}

		body := dto.CreateAuthorDTO{
			Name:         authorName,
			EnglishName:  authorName,
			OriginalName: authorName,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := authors.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/authors", h.AddAuthor)

		req := httptest.NewRequest("POST", "/authors", bytes.NewBuffer(jsonBody))
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

func AddAuthorWithTheSameNameAsAuthorOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorOnModerationID, err := moderationHelpers.CreateAuthorOnModeration(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var authorName string
		if err := env.DB.Raw("SELECT name FROM authors_on_moderation WHERE id = ?", authorOnModerationID).Scan(&authorName).Error; err != nil {
			t.Fatal(err)
		}

		body := dto.CreateAuthorDTO{
			Name:         authorName,
			EnglishName:  authorName,
			OriginalName: authorName,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := authors.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/authors", h.AddAuthor)

		req := httptest.NewRequest("POST", "/authors", bytes.NewBuffer(jsonBody))
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
