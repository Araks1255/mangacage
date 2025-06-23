package tags

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/handlers/tags"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAddTagScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                   AddTagSuccess(env),
		"with the same name as tag": AddTagWithTheSameNameAsTag(env),
		"with the same name as tag on moderation": AddTagWithTheSameNameAsTagOnModeration(env),
	}
}

func AddTagSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		body := models.TagOnModerationDTO{
			Name: uuid.New().String(),
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := tags.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/tags", h.AddTag)

		req := httptest.NewRequest("POST", "/tags", bytes.NewBuffer(jsonBody))
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

func AddTagWithTheSameNameAsTag(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		tagsIDs, err := testhelpers.CreateTags(env.DB, 1)
		if err != nil {
			t.Fatal(err)
		}

		var tagName string
		if err := env.DB.Raw("SELECT name FROM tags WHERE id = ?", tagsIDs[0]).Scan(&tagName).Error; err != nil {
			t.Fatal(err)
		}

		body := models.TagOnModerationDTO{
			Name: tagName,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := tags.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/tags", h.AddTag)

		req := httptest.NewRequest("POST", "/tags", bytes.NewBuffer(jsonBody))
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

func AddTagWithTheSameNameAsTagOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		tagOnModerationID, err := moderation.CreateTagOnModeration(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var tagOnModerationName string
		if err := env.DB.Raw("SELECT name FROM tags_on_moderation WHERE id = ?", tagOnModerationID).Scan(&tagOnModerationName).Error; err != nil {
			t.Fatal(err)
		}

		body := models.TagOnModerationDTO{
			Name: tagOnModerationName,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := tags.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/tags", h.AddTag)

		req := httptest.NewRequest("POST", "/tags", bytes.NewBuffer(jsonBody))
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
