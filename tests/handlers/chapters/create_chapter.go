package chapters

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCreateChapterScenarios(env testenv.Env) map[string]func(t *testing.T) {
	return map[string]func(t *testing.T){
		"success":                                CreateChapterSuccess(env),
		"unauthorized":                           CreateChapterByUnauthorizedUser(env),
		"non team leader":                        CreateChapterByNonTeamLeader(env),
		"does not translate title":               CreateChapterByUserWhoseTeamDoesNotTranslateTitle(env),
		"the same name as chapter on moderation": CreateChapterWithTheSameNameAsChapterOnModeration(env),
		"the same name as chapter":               CreateChapterWithTheSameNameAsChapter(env),
		"wrong volume id":                        CreateChapterWithWrongVolumeID(env),
		"invalid volume id":                      CreateChapterWithInvalidVolumeID(env),
		"wrong content type":                     CreateChapterWithWrongContentType(env),
		"without name":                           CreateChapterWithoutName(env),
		"without pages":                          CreateChapterWithoutPages(env),
	}
}

func CreateChapterSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("pages", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/test_chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := chapters.NewHandler(env.DB, env.NotificationsClient, chaptersPages)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

func CreateChapterByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		req := httptest.NewRequest("POST", "/volume/18/chapters", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func CreateChapterByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		writer := multipart.NewWriter(bytes.NewBuffer([]byte{}))
		writer.Close()

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
		req := httptest.NewRequest("POST", url, nil)
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

func CreateChapterByUserWhoseTeamDoesNotTranslateTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("pages", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/test_chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

func CreateChapterWithTheSameNameAsChapterOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("pages", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/test_chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		body2 := body

		h := chapters.NewHandler(env.DB, env.NotificationsClient, chaptersPages)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

		req2 := httptest.NewRequest("POST", url, &body2)
		req2.AddCookie(cookie)
		req2.Header.Set("Content-Type", writer.FormDataContentType())
		w2 := httptest.NewRecorder()

		r.ServeHTTP(w2, req2)

		if w2.Code != 409 {
			t.Fatal(w2.Body.String())
		}
	}
}

func CreateChapterWithTheSameNameAsChapter(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		chaptersPages := env.MongoDB.Collection(mongodb.ChaptersPagesCollection)

		h := chapters.NewHandler(env.DB, env.NotificationsClient, chaptersPages)

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

		chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, teamID, userID)
		if err != nil {
			t.Fatal(err)
		}

		var chapterName string
		env.DB.Raw("SELECT name FROM chapters WHERE id = ?", chapterID).Scan(&chapterName)
		if chapterName == "" {
			t.Fatal("ошибка при получении названия главы")
		}

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", chapterName); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("pages", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/test_chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)

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

func CreateChapterWithWrongVolumeID(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID := 9223372036854775807

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("pages", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/test_chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

func CreateChapterWithInvalidVolumeID(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		invalidVolumeID := ">_<"

		writer := multipart.NewWriter(bytes.NewBuffer([]byte{}))
		writer.Close()

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%s/chapters", invalidVolumeID)
		req := httptest.NewRequest("POST", url, nil)
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

func CreateChapterWithWrongContentType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID := 1

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

func CreateChapterWithoutName(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID := 18

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		part, err := writer.CreateFormFile("pages", "file")
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile("./test_data/test_chapter_page.png")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

func CreateChapterWithoutPages(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		volumeID := 18

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "someDescription"); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := chapters.NewHandler(env.DB, env.NotificationsClient, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader", "ex_team_leader"}))
		r.POST("/volume/:id/chapters", h.CreateChapter)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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
