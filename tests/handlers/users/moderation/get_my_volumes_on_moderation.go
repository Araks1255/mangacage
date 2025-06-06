package moderation

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	moderationHelpers "github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetMyVolumesOnModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         GetMyVolumesOnModerationSuccess(env),
		"success new":     GetMyNewVolumesOnModerationSuccess(env),
		"success edited":  GetMyEditedVolumesOnModerationSuccess(env),
		"unauthorized":    GetMyVolumesOnModerationUnauthorized(env),
		"without volumes": GetMyVolumesOnModerationWithoutVolumes(env),
		"invalid type":    GetMyVolumesOnModerationWithInvalidType(env),
		"invalid limit":   GetMyVolumesOnModerationWithInvalidLimit(env),
	}
}

func GetMyVolumesOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, userID, false); err != nil {
			t.Fatal(err)
		}

		if _, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, userID, true); err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes", nil)

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

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 2 {
			t.Fatal("не все томы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
	}
}

func GetMyNewVolumesOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, userID, false); err != nil {
				t.Fatal(err)
			}
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes?type=new", nil)

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

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 2 {
			t.Fatal("не все томы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
	}
}

func GetMyEditedVolumesOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := moderationHelpers.CreateVolumeOnModerationWithDependencies(env.DB, userID, true); err != nil {
				t.Fatal(err)
			}
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes?type=edited", nil)

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

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 2 {
			t.Fatal("не все томы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
		if _, ok := resp[0]["existing"]; !ok {
			t.Fatal("оригинальный том не дошел")
		}
		if _, ok := resp[0]["existingId"]; !ok {
			t.Fatal("id оригинального тома не дошел")
		}
	}
}

func GetMyVolumesOnModerationUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyVolumesOnModerationWithoutVolumes(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes", nil)

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

func GetMyVolumesOnModerationWithInvalidType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes?type=+_+", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf(w.Body.String())
		}
	}
}

func GetMyVolumesOnModerationWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.GET("/users/me/moderation/volumes", h.GetMyVolumesOnModeration)

		req := httptest.NewRequest("GET", "/users/me/moderation/volumes?limit=l_l", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatalf(w.Body.String())
		}
	}
}
