package moderation

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetCancelAppealForVolumeModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success new":    CancelAppealForNewVolumeModerationSuccess(env),
		"success edited": CancelAppealForEditedVolumeModerationSuccess(env),
		"other`s appeal": CancelOthersAppealForVolumeModeration(env),
	}
}

func CancelAppealForNewVolumeModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		newVolumeOnModerationID, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, userID, false)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/volumes/%d", newVolumeOnModerationID)
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

func CancelAppealForEditedVolumeModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		editedVolumeOnModerationID, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, userID, true)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/volumes/%d", editedVolumeOnModerationID)
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

func CancelOthersAppealForVolumeModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		otherUserVolumeOnModerationID, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, otherUserID, false)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/volumes/%d", otherUserVolumeOnModerationID)
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
