package users

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyProfilePictureScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":             GetMyProfilePictureSuccess(env),
		"no profile pictures": GetMyProfilePictureWithoutProfilePicture(env),
	}
}

func GetMyProfilePictureSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersProfilePictures := env.MongoDB.Collection(mongodb.UsersProfilePicturesCollection)

		data, err := os.ReadFile("./test_data/user_profile_picture.png")
		if err != nil {
			t.Fatal(err)
		}

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{ProfilePicture: data, Collection: usersProfilePictures})
		if err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, usersProfilePictures, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/profile-picture", h.GetMyProfilePicture)

		req := httptest.NewRequest("GET", "/users/me/profile-picture", nil)

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

		if len(w.Body.Bytes()) != len(data) {
			t.Fatal("возникли проблемы с аватаркой")
		}
	}
}

func GetMyProfilePictureWithoutProfilePicture(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersProfilePictures := env.MongoDB.Collection(mongodb.UsersProfilePicturesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, usersProfilePictures, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/profile-picture", h.GetMyProfilePicture)

		req := httptest.NewRequest("GET", "/users/me/profile-picture", nil)

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
