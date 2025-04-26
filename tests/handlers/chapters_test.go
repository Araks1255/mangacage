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
	"gorm.io/gorm"
)

func TestCreateChapter(t *testing.T) {
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)
	chaptersPages := env.MongoDB.Collection(constants.ChaptersPagesCollection)

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))

	r.POST("/volume/:id/chapters", h.CreateChapter)

	var err error

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

	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
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

	h := chapters.NewHandler(env.DB, chaptersOnModerationPages, chaptersPages)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))

	r.DELETE("/chapters/:id", h.DeleteChapter)

	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, userID)
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

	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, userID)
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

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	pages := make([][]byte, 1, 1)
	pages[0], err = os.ReadFile("./test_data/chapter_page.png")
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, userID, testhelpers.CreateChapterOptions{Pages: pages, Collection: chaptersPages})
	if err != nil {
		t.Fatal(err)
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

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, userID)
	if err != nil {
		t.Fatal(err)
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

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = testhelpers.CreateChapter(env.DB, volumeID, userID); err != nil {
		t.Fatal(err)
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

// Setup

func CreateVolumeWithDependencies(db *gorm.DB, userID uint) (uint, error) {
	teamID, err := testhelpers.CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	if err = testhelpers.AddUserToTeam(db, userID, teamID); err != nil {
		return 0, err
	}

	authorID, err := testhelpers.CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	titleID, err := testhelpers.CreateTitle(db, userID, authorID)
	if err != nil {
		return 0, err
	}

	if err = testhelpers.TranslateTitle(db, teamID, titleID); err != nil {
		return 0, err
	}

	volumeID, err := testhelpers.CreateVolume(db, titleID, userID)
	if err != nil {
		return 0, err
	}

	return volumeID, nil
}
