package users

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetEditProfileScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                             EditProfileSuccess(env),
		"success twice":                       EditProfileTwiceSuccess(env),
		"the same name as user":               EditProfileByAddingTheSameNameAsUser(env),
		"the same name as user on moderation": EditProfileByAddingTheSameNameAsUserOnModeration(env),
	}
}

func EditProfileSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersOnModerationProfilePictures := env.MongoDB.Collection(mongodb.UsersOnModerationProfilePicturesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, nil, usersOnModerationProfilePictures)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("userName", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("aboutYourself", "newAbout"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("profilePicture", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/user_profile_picture.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/users/me/edited", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditProfileTwiceSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		usersOnModerationProfilePictures := env.MongoDB.Collection(mongodb.UsersOnModerationProfilePicturesCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderation.CreateUserOnModeration(env.DB, moderation.CreateUserOnModerationOptions{ExistingID: userID}); err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, nil, usersOnModerationProfilePictures)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("userName", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("aboutYourself", "newAbout"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("profilePicture", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/user_profile_picture.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/users/me/edited", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditProfileUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := users.NewHandler(env.DB, env.NotificationsClient, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		req := httptest.NewRequest("POST", "/users/me/edited", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func EditProfileWithoutEditableParams(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("randomField", "L_L"); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/users/me/edited", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditProfileWithTooLargeProfilePicture(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		part, err := writer.CreateFormFile("profilePicture", "file")
		if err != nil {
			t.Fatal(err)
		}
		data := make([]byte, 3<<20, 3<<20)
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/users/me/edited", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditProfileByAddingTheSameNameAsUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		existingUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var existingUserName string
		if err := env.DB.Raw("SELECT user_name FROM users WHERE id = ?", existingUserID).Scan(&existingUserName).Error; err != nil {
			t.Fatal(err)
		}

		if existingUserName == "" {
			t.Fatal("не удалось получить имя пользователя")
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("userName", existingUserName); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/users/me/edited", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditProfileByAddingTheSameNameAsUserOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		userOnModerationID, err := moderation.CreateUserOnModeration(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var userOnModerationName string
		if err := env.DB.Raw("SELECT user_name FROM users_on_moderation WHERE id = ?", userOnModerationID).Scan(&userOnModerationName).Error; err != nil {
			t.Fatal(err)
		}

		if userOnModerationName == "" {
			t.Fatal("не удалось получить имя пользователя")
		}

		h := users.NewHandler(env.DB, env.NotificationsClient, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/users/me/edited", h.EditProfile)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("userName", userOnModerationName); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/users/me/edited", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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
