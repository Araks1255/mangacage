package moderation

import (
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

func GetCancelAppealForProfileChangesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         CancelAppealForProfileChangesSuccess(env),
		"without changes": CancelAppealForProfileChangesWithoutProfileChanges(env),
		"unauthorized":    CancelAppealForProfileChangesUnauthorized(env),
	}
}

func CancelAppealForProfileChangesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersProfilePictures := env.MongoDB.Collection(mongodb.UsersOnModerationProfilePicturesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile("./test_data/user_profile_picture.png")
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateUserOnModeration(
			env.DB, moderationHelpers.CreateUserOnModerationOptions{ExistingID: userID, ProfilePicture: data, Collection: usersProfilePictures},
		); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, usersProfilePictures, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/profile/edited", h.CancelAppealForProfileChanges)

		req := httptest.NewRequest("DELETE", "/users/me/moderation/profile/edited", nil)

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

func CancelAppealForProfileChangesWithoutProfileChanges(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/profile/edited", h.CancelAppealForProfileChanges)

		req := httptest.NewRequest("DELETE", "/users/me/moderation/profile/edited", nil)

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

func CancelAppealForProfileChangesUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/profile/edited", h.CancelAppealForProfileChanges)

		req := httptest.NewRequest("DELETE", "/users/me/moderation/profile/edited", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}
