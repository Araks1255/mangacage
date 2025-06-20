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

func GetCreateTeamScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                             CreateTeamSuccess(env),
		"unauthorized":                        CreateTeamUnauthorizedUser(env),
		"user already in team":                CreateChapterByUserThatAlreadyInTeam(env),
		"user already has team on moderation": CreateTeamByUserThatAlredyHasTeamOnModeration(env),
		"the same name as team":               CreateTeamWithTheSameNameAsTeam(env),
		"the same name as team on moderation": CreateTeamWithTheSameNameAsTeamOnModeration(env),
		"too large cover":                     CreateTeamWithTooLargeCover(env),
		"without name":                        CreateTeamWithoutName(env),
		"without cover":                       CreateTeamWithoutCover(env),
	}
}

func CreateTeamSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "teamDescription"); err != nil {
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

		h := teams.NewHandler(env.DB, teamsOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func CreateChapterByUserThatAlreadyInTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
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

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "teamDescription"); err != nil {
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

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamByUserThatAlredyHasTeamOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderation.CreateTeamOnModeration(env.DB, userID); err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "teamDescription"); err != nil {
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

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamWithTheSameNameAsTeamOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
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
			t.Fatal("не удалось получить название созданной главы на модерации")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", teamOnModerationName); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "teamDescription"); err != nil {
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

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamWithTheSameNameAsTeam(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		existingTeamID, err := testhelpers.CreateTeam(env.DB, userID)
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
		if err = writer.WriteField("description", "teamDescription"); err != nil {
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

		h := teams.NewHandler(env.DB, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamWithoutName(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("description", "teamDescription"); err != nil {
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

		h := teams.NewHandler(env.DB, teamsOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamWithoutCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "teamDescription"); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, teamsOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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

func CreateTeamWithTooLargeCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "teamDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		cover := make([]byte, 3<<20)
		if _, err = part.Write(cover); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := teams.NewHandler(env.DB, teamsOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/teams", h.CreateTeam)

		req := httptest.NewRequest("POST", "/teams", &body)
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
