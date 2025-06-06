package moderation

import (
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyProfileChangesOnModerationScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":         GetMyProfileChangesOnModerationSuccess(env),
		"unauthorized":    GetMyProfileChangesOnModerationUnauthorized(env),
		"without changes": GetMyProfileChangesOnModerationWithoutChanges(env),
	}
}

func GetMyProfileChangesOnModerationSuccess(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateUserOnModeration(
			env.DB, moderationHelpers.CreateUserOnModerationOptions{ExistingID: userID},
		); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/profile", h.GetMyProfileChangesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile", nil)

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

func GetMyProfileChangesOnModerationUnauthorized(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/profile", h.GetMyProfileChangesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyProfileChangesOnModerationWithoutChanges(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/profile", h.GetMyProfileChangesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile", nil)

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
