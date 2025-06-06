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
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCreateTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                              CreateTitleSuccess(env),
		"unauthorized":                         CreateTitleUnauthorized(env),
		"without name":                         CreateTitleWithoutName(env),
		"without author":                       CreateTitleWithoutAuthor(env),
		"without genres":                       CreateTitleWithoutGenres(env),
		"without cover":                        CreateTitleWithoutCover(env),
		"too large cover":                      CreateTitleWithTooLargeCover(env),
		"invalid author id":                    CreateTitleWithInvalidAuthorId(env),
		"invalid genres ids":                   CreateTitleWithInvalidGenres(env),
		"wrong author id":                      CreateTitleWithWrongAuthorAuthor(env),
		"wrong genres ids":                     CreateTitleWithWrongGenres(env),
		"the same name as title":               CreateTitleWithTheSameNameAsTitle(env),
		"the same name as title on moderation": CreateTitleWithTheSameNameAsTitleOnModerationOnModeration(env),
	}
}

func CreateTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func CreateTitleWithoutName(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genreId", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithoutAuthor(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genreId", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithoutGenres(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithoutCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genreId", genresIDs[i]); err != nil {
				t.Fatal(err)
			}
		}

		writer.Close()

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithTooLargeCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
				t.Fatal(err)
			}
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		data := make([]byte, 3<<20, 3<<20)
		if _, err := part.Write(data); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithInvalidAuthorId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidAuthorID := "U_U"

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", invalidAuthorID); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithInvalidGenres(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidGenresIDs := []string{"<_>", ">_<"}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(invalidGenresIDs); i++ {
			if err := writer.WriteField("genresIds", invalidGenresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithWrongAuthorAuthor(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID := 9223372036854775807

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithWrongGenres(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs := []string{"9223372036854775806", "9223372036854775805"}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithTheSameNameAsTitle(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
		}

		titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID)
		if err != nil {
			t.Fatal(err)
		}

		var titleName string
		if err := env.DB.Raw("SELECT name FROM titles WHERE id = ?", titleID).Scan(&titleName).Error; err != nil {
			t.Fatal(err)
		}

		if titleName == "" {
			t.Fatal("не удалось получить название тайтла")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", titleName); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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

func CreateTitleWithTheSameNameAsTitleOnModerationOnModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		var genresIDs []string
		env.DB.Raw("SELECT id FROM genres LIMIT 2").Scan(&genresIDs)
		if len(genresIDs) == 0 {
			t.Fatal("не удалось получить жанры")
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
			t.Fatal("не удалось получить название тайтла на модерации")
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", titleOnModerationName); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}
		if err := writer.WriteField("description", "desc"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", genresIDs[i]); err != nil {
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

		h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", &body)
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
