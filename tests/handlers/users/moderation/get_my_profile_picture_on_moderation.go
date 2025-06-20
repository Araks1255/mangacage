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

func GetGetMyProfilePictureOnModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                 GetMyProfilePictureOnModerationSuccess(env),
		"unauthorized":            GetMyProfilePictureOnModerationUnauthorized(env),
		"without profile picture": GetMyProfilePictureOnModerationWithoutProfilePicture(env),
		"without changes":         GetMyProfilePictureOnModerationWithoutProfileChanges(env),
	}
}

func GetMyProfilePictureOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersProfilePictures := env.MongoDB.Collection(mongodb.UsersProfilePicturesCollection)

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
		r.GET("/users/me/moderation/profile/picture", h.GetMyProfilePictureOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile/picture", nil)

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

func GetMyProfilePictureOnModerationUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/profile/picture", h.GetMyProfilePictureOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile/picture", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyProfilePictureOnModerationWithoutProfilePicture(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersProfilePictures := env.MongoDB.Collection(mongodb.UsersProfilePicturesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateUserOnModeration(
			env.DB, moderationHelpers.CreateUserOnModerationOptions{ExistingID: userID, Collection: usersProfilePictures},
		); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, usersProfilePictures, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/profile/picture", h.GetMyProfilePictureOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile/picture", nil)

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

func GetMyProfilePictureOnModerationWithoutProfileChanges(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersProfilePictures := env.MongoDB.Collection(mongodb.UsersProfilePicturesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, usersProfilePictures, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/profile/picture", h.GetMyProfilePictureOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/profile/picture", nil)

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
