package moderation

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetCancelAppealForTeamModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success new":    CancelAppealForNewTeamModerationSuccess(env),
		"success edited": CancelAppealForEditedTeamModerationSuccess(env),
		"other`s appeal": CancelOthersAppealForTeamModeration(env),
	}
}

func CancelAppealForNewTeamModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile("./test_data/team_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		newTeamOnModerationID, err := moderationHelpers.CreateTeamOnModeration(
			env.DB, userID, moderationHelpers.CreateTeamOnModerationOptions{Cover: data, Collection: teamsCovers},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, teamsCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/teams/%d", newTeamOnModerationID)
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

func CancelAppealForEditedTeamModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		existingTeamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile("./test_data/team_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		editedTeamOnModerationID, err := moderationHelpers.CreateTeamOnModeration(
			env.DB, userID, moderationHelpers.CreateTeamOnModerationOptions{Cover: data, Collection: teamsCovers, ExistingID: existingTeamID},
		)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, teamsCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/teams/%d", editedTeamOnModerationID)
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

func CancelOthersAppealForTeamModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserTeamOnModerationID, err := moderationHelpers.CreateTeamOnModeration(env.DB, otherUserID)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, teamsCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/teams/%d", otherUserTeamOnModerationID)
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
