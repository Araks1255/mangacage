package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/internal/testhelpers"
	"github.com/Araks1255/mangacage/pkg/constants"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func TestCreateChapter(t *testing.T) {
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))

	r.POST("/volume/:id/chapters", h.CreateChapter)

	var userID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&userID)
	if userID == 0 {
		t.Fatal("Юзер не найден")
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err = writer.WriteField("name", "testChapter"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("description", "someDescription"); err != nil {
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

	var volumeID uint
	env.DB.Raw("SELECT id FROM volumes WHERE name = 'volume_test'").Scan(&volumeID)
	if volumeID == 0 {
		t.Fatal("Том не найден")
	}

	url := fmt.Sprintf("/volume/%d/chapters", volumeID)

	req := httptest.NewRequest("POST", url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.AddCookie(&http.Cookie{
		Name:  "mangacage_token",
		Value: tokenString,
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatal(w.Body.String())
	}
}

func TestDeleteChapter(t *testing.T) {
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)

	var userID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&userID)
	if userID == 0 {
		t.Fatal("Юзер не найден")
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))

	r.DELETE("/chapters/:id", h.DeleteChapter)

	chapterID, err := testhelpers.CreateTestChapter(env.DB, chaptersPages)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/chapters/%d", chapterID)

	req := httptest.NewRequest("DELETE", url, nil)

	req.AddCookie(&http.Cookie{
		Name:  "mangacage_token",
		Value: tokenString,
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestEditChapter(t *testing.T) {
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)

	var userID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&userID)
	if userID == 0 {
		t.Fatal("Юзер не найден")
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))

	r.POST("/chapters/:id/edited", h.EditChapter)

	body := map[string]any{
		"name":        "chapterTest",
		"description": "someDescription",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateTestChapter(env.DB, chaptersPages)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/chapters/%d/edited", chapterID)

	req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{
		Name:  "mangacage_token",
		Value: tokenString,
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatal(w.Body.String())
	}
}

func TestGetChapterPage(t *testing.T) {
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)

	var chapterID uint
	env.DB.Raw("SELECT id FROM chapters WHERE name = 'chapter_test'").Scan(&chapterID)
	if chapterID == 0 {
		t.Fatal("Тестовая глава не найдена")
	}

	h := chapters.NewHandler(env.DB, chaptersPages, chaptersOnModerationPages)

	r := gin.New()
	r.GET("/chapters/:id/page/:page", h.GetChapterPage)

	url := fmt.Sprintf("/chapters/%d/page/0", chapterID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestGetChapter(t *testing.T) {
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)

	var chapterID uint
	env.DB.Raw("SELECT id FROM chapters WHERE name = 'chapter_test'").Scan(&chapterID)
	if chapterID == 0 {
		t.Fatal("Тестовая глава не найдена")
	}

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()

	r.GET("/chapters/:id", h.GetChapter)

	url := fmt.Sprintf("/chapters/%d", chapterID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestGetVolumeChapters(t *testing.T) {
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)

	var volumeID uint
	env.DB.Raw("SELECT volume_id FROM chapters WHERE name = 'chapter_test'").Scan(&volumeID)
	if volumeID == 0 {
		t.Fatal("Тестовый том не найден")
	}

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()
	r.GET("/volume/:id/chapters", h.GetVolumeChapters)

	url := fmt.Sprintf("/volume/%d/chapters", volumeID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}
