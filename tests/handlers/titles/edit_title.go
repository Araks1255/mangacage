package titles

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	titlesHelpers "github.com/Araks1255/mangacage/testhelpers/titles"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetEditTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                              EditTitleSuccess(env),
		"success twice":                        EditTitleTwiceSuccess(env),
		"unauthorized":                         EditTitleByUnauthorizedUser(env),
		"non team leader":                      EditTitleByNonTeamLeader(env),
		"without parameters":                   EditTitleWithoutEditableParams(env),
		"invalid title id":                     EditTitleWithInvalidTitleId(env),
		"invalid author id":                    EditTitleWithInvalidAuthorId(env),
		"invalid genres ids":                   EditTitleWithInvalidGenresIds(env),
		"wrong title id":                       EditTitleWithWrongTitleId(env),
		"wrong author id":                      EditTitleWithWrongAuthorId(env),
		"wrong genres ids":                     EditTitleWithWrongGenresIds(env),
		"the same name as title":               EditTitleByAddingTheSameNameAsTitle(env),
		"the same name as title on moderation": EditTitleByAddingTheSameNameAsTitleOnModeration(env),
		"user team does not translate title":   EditTitleByUserWhoseTeamDoesNotTranslateTitle(env),
		"wrong type":                           EditTitleWithWrongType(env),
		"wrong publishing status":              EditTitleWithWrongPublishingStatus(env),
	}
}

func EditTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs, err := testhelpers.CreateGenres(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		tagsIDs, err := testhelpers.CreateTags(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		cover, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		name := uuid.New().String()
		description := "newDescription"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, nil, nil, nil, nil, &description, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", contentType)

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

func EditTitleTwiceSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := moderation.CreateTitleOnModeration(env.DB, userID); err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []uint
		if err := env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs).Error; err != nil {
			t.Fatal(err)
		}

		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "newDescription"); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err = writer.WriteField("genresIds", fmt.Sprintf("%d", genresIDs[i])); err != nil {
				t.Fatal(err)
			}
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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

func EditTitleByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		req := httptest.NewRequest("POST", "/titles/18/edited", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func EditTitleByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		req := httptest.NewRequest("POST", "/titles/18/edited", nil)

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

func EditTitleWithoutEditableParams(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", contentType)

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

func EditTitleWithInvalidTitleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidTitleID := "o_o"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		url := fmt.Sprintf("/titles/%s/edited", invalidTitleID)
		req := httptest.NewRequest("POST", url, body)
		req.Header.Set("Content-Type", contentType)

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

func EditTitleWithInvalidAuthorId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidAuthorID := "O_o"

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("authorId", invalidAuthorID); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/titles/18/edited", &body)
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

func EditTitleWithInvalidGenresIds(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		invalidGenresIDs := []string{"c_c", "'-'"}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		for i := 0; i < len(invalidGenresIDs); i++ {
			if err = writer.WriteField("genresIds", invalidGenresIDs[i]); err != nil {
				t.Fatal(err)
			}
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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

func EditTitleWithTooLargeCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		data := make([]byte, 3<<20, 3<<20)
		if _, err := part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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

func EditTitleByUserWhoseTeamDoesNotTranslateTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditTitleByAddingTheSameNameAsTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		existingTitleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var existingTitleName string
		if err := env.DB.Raw("SELECT name FROM titles WHERE id = ?", existingTitleID).Scan(&existingTitleName).Error; err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", existingTitleName); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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

func EditTitleByAddingTheSameNameAsTitleOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		titleOnModerationID, err := moderation.CreateTitleOnModeration(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var titleOnModerationName string
		if err := env.DB.Raw("SELECT name FROM titles_on_moderation WHERE id = ?", titleOnModerationID).Scan(&titleOnModerationName).Error; err != nil {
			t.Fatal(err)
		}

		if titleOnModerationName == "" {
			t.Fatal("не удалось получить id тайтла на модерации")
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", titleOnModerationName); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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

func EditTitleWithWrongTitleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID := 9223372036854775807

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditTitleWithWrongAuthorId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		authorID := 9223372036854775807

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err = writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("description", "newDescription"); err != nil {
			t.Fatal(err)
		}
		if err = writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditTitleWithWrongGenresIds(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs := []uint{9223372036854775807, 9223372036854775806}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		for i := 0; i < len(genresIDs); i++ {
			if err = writer.WriteField("genresIds", fmt.Sprintf("%d", genresIDs[i])); err != nil {
				t.Fatal(err)
			}
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func EditTitleWithWrongType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("type", "M_M"); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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

func EditTitleWithWrongPublishingStatus(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/titles/:id/edited", h.EditTitle)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("publishingStatus", "F_F"); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/titles/%d/edited", titleID)
		req := httptest.NewRequest("POST", url, &body)
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
