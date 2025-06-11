package volumes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/volumes"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCreateVolumeScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                                     CreateVolumeSuccess(env),
		"unauthorized":                                CreateVolumeUnauthorized(env),
		"by non team leader":                          CreateVolumeByNonTeamLeader(env),
		"in title translating by other team":          CreateVolumeInTitleTranslatingByOtherTeam(env),
		"the same name as volume":                     CreateVolumeWithSameNameAsVolume(env),
		"the same name as volume on moderation":       CreateVolumeWithSameNameAsVolumeOnModeration(env),
		"wrong title id":                              CreateVolumeWithWrongTitleId(env),
		"invalid title id":                            CreateVolumeWithInvalidTitleId(env),
		"without name":                                CreateVolumeWithoutName(env),
		"by user whose team does not translate title": CreateVolumeByUserWhoseTeamDoesNotTranslateTitle(env),
	}
}

func CreateVolumeSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatalf(w.Body.String())
		}
	}
}

func CreateVolumeUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		req := httptest.NewRequest("POST", "/titles/18/volumes", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func CreateVolumeByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		req := httptest.NewRequest("POST", "/titles/18/volumes", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatalf(w.Body.String())
		}
	}
}

func CreateVolumeByUserWhoseTeamDoesNotTranslateTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func CreateVolumeInTitleTranslatingByOtherTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		otherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, otherUserID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func CreateVolumeWithSameNameAsVolume(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.TranslateTitle(env.DB, teamID, titleID); err != nil {
			t.Fatal(err)
		}

		existingVolumeID, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		var existingVolumeName string
		if err := env.DB.Raw("SELECT name FROM volumes WHERE id = ?", existingVolumeID).Scan(&existingVolumeName).Error; err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        existingVolumeName,
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 409 {
			t.Fatalf(w.Body.String())
		}
	}
}

func CreateVolumeWithSameNameAsVolumeOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		existingVolumeOnModerationID, err := moderation.CreateVolumeOnModeration(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		var existingVolumeOnModerationName string

		if err := env.DB.Raw(
			"SELECT name FROM volumes_on_moderation WHERE id = ?", existingVolumeOnModerationID,
		).Scan(&existingVolumeOnModerationName).Error; err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        existingVolumeOnModerationName,
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 409 {
			t.Fatalf(w.Body.String())
		}
	}
}

func CreateVolumeWithWrongTitleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID := 9223372036854775807

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func CreateVolumeWithInvalidTitleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/titles/$_$/volumes", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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

func CreateVolumeWithoutName(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/titles/:id/volumes", h.CreateVolume)

		body := map[string]any{
			"description": "some description",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

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
