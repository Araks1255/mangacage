package titles

import (
	"bytes"
	"fmt"
	"log"
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

func GetCreateTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":                              CreateTitleSuccess(env),
		"with author on moderation success":    CreateTitleWithAuthorOnModerationSuccess(env),
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
		"wrong publishing status":              CreateTitleWithWrongPublishingStatus(env),
		"wrong type":                           CreateTitleWithWrongType(env),
		"the same name as title":               CreateTitleWithTheSameNameAsTitle(env),
		"the same name as title on moderation": CreateTitleWithTheSameNameAsTitleOnModerationOnModeration(env),
		"with two authors":                     CreateTitleWithTwoAuthors(env),
	}
}

func CreateTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleWithAuthorOnModerationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorOnModerationID, err := moderation.CreateAuthorOnModeration(env.DB, userID)
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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, nil, &authorOnModerationID, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, env.NotificationsClient, nil)

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
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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

		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			nil, nil, nil, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleWithoutAuthor(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, nil, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleWithoutGenres(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, nil, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleWithoutCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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

		name := uuid.New().String()
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleWithTooLargeCover(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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

		cover := make([]byte, 3<<20, 3<<20)

		name := uuid.New().String()
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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

func CreateTitleWithInvalidAuthorId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("englishName", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("originalName", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("ageLimit", "18"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("type", "manga"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("publishingStatus", "ongoing"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("yearOfRelease", "1999"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("authorId", "R_R"); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(genresIDs); i++ {
			if err := writer.WriteField("genresIds", fmt.Sprintf("%d", genresIDs[i])); err != nil {
				t.Fatal(err)
			}
		}

		for i := 0; i < len(tagsIDs); i++ {
			if err := writer.WriteField("tagsIds", fmt.Sprintf("%d", tagsIDs[i])); err != nil {
				t.Fatal(err)
			}
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(cover); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

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
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidGenresIDs := []string{"<_>", ">_<"}

		tagsIDs, err := testhelpers.CreateTags(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		cover, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)

		if err := writer.WriteField("name", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("englishName", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("originalName", uuid.New().String()); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("ageLimit", "18"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("type", "manga"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("publishingStatus", "ongoing"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("yearOfRelease", "1999"); err != nil {
			t.Fatal(err)
		}

		if err := writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < len(invalidGenresIDs); i++ {
			if err := writer.WriteField("genresIds", invalidGenresIDs[i]); err != nil {
				t.Fatal(err)
			}
		}

		for i := 0; i < len(tagsIDs); i++ {
			if err := writer.WriteField("tagsIds", fmt.Sprintf("%d", tagsIDs[i])); err != nil {
				t.Fatal(err)
			}
		}

		part, err := writer.CreateFormFile("cover", "file")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(cover); err != nil {
			t.Fatal(err)
		}

		writer.Close()

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

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
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID := uint(9223372036854775807)

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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
		req.Header.Set("Content-Type", contentType)

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
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		genresIDs := []uint{9223372036854775807, 9223372036854775806}

		tagsIDs, err := testhelpers.CreateTags(env.DB, 2)
		if err != nil {
			t.Fatal(err)
		}

		cover, err := os.ReadFile("./test_data/title_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		name := uuid.New().String()
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
		req.Header.Set("Content-Type", contentType)

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
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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

		existingTitleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var name string
		if err := env.DB.Raw("SELECT name FROM titles WHERE id = ?", existingTitleID).Scan(&name).Error; err != nil {
			t.Fatal(err)
		}
		log.Println(name, "<===")

		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
		req.Header.Set("Content-Type", contentType)

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
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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

		titleOnModerationID, err := moderation.CreateTitleOnModeration(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var name string
		if err := env.DB.Raw("SELECT name FROM titles_on_moderation WHERE id = ?", titleOnModerationID).Scan(&name).Error; err != nil {
			t.Fatal(err)
		}

		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
		req.Header.Set("Content-Type", contentType)

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

func CreateTitleWithWrongPublishingStatus(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "._."
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
		req.Header.Set("Content-Type", contentType)

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

func CreateTitleWithWrongType(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
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
		ageLimit := "18"
		titleType := "._."
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, nil, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
		req.Header.Set("Content-Type", contentType)

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

func CreateTitleWithTwoAuthors(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorOnModerationID, err := moderation.CreateAuthorOnModeration(env.DB, userID)
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
		ageLimit := "18"
		titleType := "manga"
		publishingStatus := "ongoing"
		yearOfRelease := "1999"

		body, contentType, err := titlesHelpers.FillTitleRequestBody(
			&name, &name, &name, &ageLimit, &titleType, &publishingStatus, &yearOfRelease, nil, &authorID, &authorOnModerationID, genresIDs, tagsIDs, cover,
		)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, env.NotificationsClient, titlesOnModerationCovers)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles", h.CreateTitle)

		req := httptest.NewRequest("POST", "/titles", body)
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
