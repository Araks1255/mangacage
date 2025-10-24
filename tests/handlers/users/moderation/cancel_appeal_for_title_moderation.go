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

func GetCancelAppealForTitleModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success new":    CancelAppealForNewTitleModerationSuccess(env),
		"success edited": CancelAppealForEditedTitleModerationSuccess(env),
		"other`s appeal": CancelOthersAppealForTitleModeration(env),
	}
}

func CancelAppealForNewTitleModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		newTitleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(
			env.DB, userID, moderationHelpers.CreateTitleOnModerationOptions{Cover: data, Collection: titlesCovers, Genres: []string{"action", "fighting"}},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, titlesCovers, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/titles/%d", newTitleOnModerationID)
		req := httptest.NewRequest("DELETE", url, nil)

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
	}
}

func CancelAppealForEditedTitleModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		existingTitleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		editedTitleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(
			env.DB, userID,
			moderationHelpers.CreateTitleOnModerationOptions{Cover: data, Collection: titlesCovers, Genres: []string{"action", "fighting"}, ExistingID: existingTitleID},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, titlesCovers, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/titles/%d", editedTitleOnModerationID)
		req := httptest.NewRequest("DELETE", url, nil)

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
	}
}

func CancelOthersAppealForTitleModeration(env testenv.Env) func(*testing.T) {
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
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/titles/%d", otherUserTitleOnModerationID)
		req := httptest.NewRequest("DELETE", url, nil)

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
