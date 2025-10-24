package moderation

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyTitleOnModerationCoverScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":       GetMyTitleOnModerationCoverSuccess(env),
		"unauthorized":  GetMyTitleOnModerationCoverUnauthorized(env),
		"other`s title": GetOthersTitleOnModerationCover(env),
		"without cover": GetMyTitleOnModerationCoverWithoutCover(env),
		"wrong id":      GetMyTitleOnModerationCoverWithWrongId(env),
		"invalid id":    GetMyTitleOnModerationCoverWithInvalidId(env),
	}
}

func GetMyTitleOnModerationCoverSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		cover, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		titleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(
			env.DB, userID, moderationHelpers.CreateTitleOnModerationOptions{Cover: cover, Collection: titlesCovers},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, titlesCovers, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

		url := fmt.Sprintf("/users/me/moderation/titles/%d/cover", titleOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

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

		if len(w.Body.Bytes()) != len(cover) {
			t.Fatal("возникли проблемы с получением обложки")
		}
	}
}

func GetMyTitleOnModerationCoverUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles/18/cover", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetOthersTitleOnModerationCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserTitleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(env.DB, otherUserID)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, titlesCovers, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

		url := fmt.Sprintf("/users/me/moderation/titles/%d/cover", otherUserTitleOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

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

func GetMyTitleOnModerationCoverWithoutCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, titlesCovers, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

		url := fmt.Sprintf("/users/me/moderation/titles/%d/cover", titleOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

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

func GetMyTitleOnModerationCoverWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleOnModerationID := 9223372036854775807

		h := moderation.NewHandler(env.DB, titlesCovers, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

		url := fmt.Sprintf("/users/me/moderation/titles/%d/cover", titleOnModerationID)
		req := httptest.NewRequest("GET", url, nil)

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

func GetMyTitleOnModerationCoverWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

		req := httptest.NewRequest("GET", "/users/me/moderation/titles/Ъ_Ъ/cover", nil)

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
