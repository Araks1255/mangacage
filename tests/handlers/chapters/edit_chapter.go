package chapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetEditChapterScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":                                EditChapterSuccess(env),
		"unauthorized":                           EditChapterByUnauthorizedUser(env),
		"non team leader":                        EditChapterByNonTeamLeader(env),
		"user team does not translate the title": EditChapterByUserWhoseTeamDoesNotTranslateTitle(env),
		"wrong volume id":                        EditChapterWithWrongVolumeId(env),
		"the same name as chapter on moderation": EditChapterByAddingTheSameNameAsChapterOnModeration(env),
		"the same name as chapter":               EditChapterByAddingTheSameNameAsChapter(env),
		"invalid chapter id":                     EditChapterWithInvalidChapterId(env),
		"without editable parameters":            EditChapterWithoutEditableParameters(env),
		"wrong content type":                     EditChapterWithWrongContentType(env),
	}
}

func EditChapterSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

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

		chapterID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "someDescription",
			"volume":      1,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, chaptersPages)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%d/edited", chapterID)
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
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		req := httptest.NewRequest("POST", "/chapters/18/edited", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		req := httptest.NewRequest("POST", "/chapters/18/edited", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterByUserWhoseTeamDoesNotTranslateTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		body := map[string]any{
			"description": "someDescription",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%d/edited", chapterID)
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
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterWithWrongVolumeId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		body := map[string]any{
			"name":        uuid.New().String(),
			"description": "someDescription",
			"volumeId":    9223372036854775807,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%d/edited", chapterID)
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
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterByAddingTheSameNameAsChapterOnModeration(env testenv.Env) func(*testing.T) {
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

		existingChapter1ID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		existingChapter2ID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		body := map[string]any{
			"name":        "the same name",
			"description": "someDescription",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%d/edited", existingChapter1ID)
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
			t.Fatal(w.Body.String())
		}

		url2 := fmt.Sprintf("/chapters/%d/edited", existingChapter2ID)
		req2 := httptest.NewRequest("POST", url2, bytes.NewBuffer(jsonBody))
		req2.AddCookie(cookie)

		w2 := httptest.NewRecorder()

		r.ServeHTTP(w2, req2)

		if w2.Code != 409 {
			t.Fatal(w2.Body.String())
		}
	}
}

func EditChapterByAddingTheSameNameAsChapter(env testenv.Env) func(*testing.T) {
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

		existingChapterID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		var existingChapterName string
		env.DB.Raw("SELECT name FROM chapters WHERE id = ?", existingChapterID).Scan(&existingChapterName)
		if existingChapterName == "" {
			t.Fatal("ошибка при получении названия существующей главы")
		}

		editableChapterID, err := testhelpers.CreateChapter(env.DB, titleID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		body := map[string]any{
			"name":        existingChapterName,
			"description": "someDescription",
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%d/edited", editableChapterID)
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
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterWithInvalidChapterId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidChapterID := "^_^"

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%s/edited", invalidChapterID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", "application/json")

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

func EditChapterWithoutEditableParameters(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		chapterID := 18

		body := map[string]any{
			"random parameter": 18,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		url := fmt.Sprintf("/chapters/%d/edited", chapterID)
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
			t.Fatal(w.Body.String())
		}
	}
}

func EditChapterWithWrongContentType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/chapters/:id/edited", h.EditChapter)

		req := httptest.NewRequest("POST", "/chapters/18/edited", nil)
		req.Header.Set("Content-Type", "text/plain")

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
