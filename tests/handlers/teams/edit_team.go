package teams

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetEditTeamScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                             EditTeamSuccess(env),
		"unauthorized":                        EditTeamByUnauthorizedUser(env),
		"non team leader":                     EditTeamByNonTeamLeader(env),
		"no parameters":                       EditTeamWithNoEditableParameters(env),
		"too large cover":                     EditTeamWithTooLargeCover(env),
		"the same name as team":               EditTeamByAddingTheSameNameAsTeam(env),
		"the same name as team on moderation": EditTeamByAddingTheSameNameAsTeamOnModeration(env),
	}
}

func EditTeamSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "newDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/team_cover.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, teamsOnModerationCovers, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", &body)
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

func EditTeamByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := teams.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func EditTeamByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := teams.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", nil)

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

func EditTeamWithNoEditableParameters(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("randomParameter", ";-("); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", &body)
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

func EditTeamWithTooLargeCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		cover := make([]byte, 3<<20)
		if _, err = part.Write(cover); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", &body)
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

func EditTeamByAddingTheSameNameAsTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		anotherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		existingTeamID, err := testhelpers.CreateTeam(env.DB, anotherUserID)
		if err != nil {
			t.Fatal(err)
		}

		var existingTeamName string
		env.DB.Raw("SELECT name FROM teams WHERE id = ?", existingTeamID).Scan(&existingTeamName)
		if existingTeamName == "" {
			t.Fatal("не удалось получить название созданной команды")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", existingTeamName); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", &body)
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

func EditTeamByAddingTheSameNameAsTeamOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err = testhelpers.AddUserToTeam(env.DB, userID, teamID); err != nil {
			t.Fatal(err)
		}

		anotherUserID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamOnModerationID, err := moderation.CreateTeamOnModeration(env.DB, anotherUserID)
		if err != nil {
			t.Fatal(err)
		}

		var teamOnModerationName string
		env.DB.Raw("SELECT name FROM teams_on_moderation WHERE id = ?", teamOnModerationID).Scan(&teamOnModerationName)
		if teamOnModerationName == "" {
			t.Fatal("не удалось получить название созданной команды на модерации")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", teamOnModerationName); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "newDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/team_cover.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.POST("/teams/my/edited", h.EditTeam)

		req := httptest.NewRequest("POST", "/teams/my/edited", &body)
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
