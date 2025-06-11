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

func GetEditVolumeScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":                               EditVolumeSuccess(env),
		"unauthorized":                          EditVolumeUnauthorized(env),
		"user team does not translate title":    EditVolumeByUserWhoseTeamDoesNotTranslateTitle(env),
		"the same name as volume on moderation": EditVolumeWithTheSameNameAsVolumeOnModeration(env),
		"wrong id":                              EditVolumeWithWrongId(env),
		"invalid id":                            EditVolumeWithInvalidId(env),
		"by non team leader":                    EditVolumeByNonTeamLeader(env),
		"with the same name as volume":          EditVolumeWithTheSameNameAsVolume(env),
		"without editable parameters":           EditVolumeWithoutEditableParameters(env),
	}
}

func EditVolumeSuccess(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "someDescription",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/volumes/%d/edited", volumeID)
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

func EditVolumeUnauthorized(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		req := httptest.NewRequest("POST", "/volumes/18/edited", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatalf(w.Body.String())
		}
	}
}

func EditVolumeByUserWhoseTeamDoesNotTranslateTitle(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "desc",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/volumes/%d/edited", volumeID)
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

func EditVolumeByNonTeamLeader(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/volumes/:id/edited", h.EditVolume)

		req := httptest.NewRequest("POST", "/volumes/18/edited", nil)
		req.Header.Set("Content-Type", "application/json")

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

func EditVolumeWithTheSameNameAsVolume(env testenv.Env) func(t *testing.T) {
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

		volumeID, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID)
		if err != nil {
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
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"name":        existingVolumeName,
			"description": "desc",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/volumes/%d/edited", volumeID)
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

func EditVolumeWithTheSameNameAsVolumeOnModeration(env testenv.Env) func(t *testing.T) {
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

		volumeID, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		volumeOnModerationID, err := moderation.CreateVolumeOnModeration(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		var volumeOnModerationName string
		if err := env.DB.Raw("SELECT name FROM volumes_on_moderation WHERE id = ?", volumeOnModerationID).Scan(&volumeOnModerationName).Error; err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"name":        volumeOnModerationName,
			"description": "desc",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/volumes/%d/edited", volumeID)
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

func EditVolumeWithWrongId(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		wrongID := 9223372036854775807

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "desc",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/volumes/%d/edited", wrongID)
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

func EditVolumeWithInvalidId(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "desc",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest("POST", "/volumes/._./edited", bytes.NewBuffer(jsonBody))
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

func EditVolumeWithoutEditableParameters(env testenv.Env) func(t *testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/volumes/:id/edited", h.EditVolume)

		body := map[string]any{
			"random": "field",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("/volumes/%d/edited", volumeID)
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
