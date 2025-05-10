package handlers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/internal/testhelpers"
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func TestCreateTitle(t *testing.T) {
	titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

	creatorID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(creatorID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	var genres []string
	env.DB.Raw("SELECT name FROM genres").Scan(&genres)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err = writer.WriteField("name", "title"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("description", "someDescription"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("authorId", fmt.Sprintf("%d", authorID)); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(genres); i++ {
		if err = writer.WriteField("genres", genres[i]); err != nil {
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

func TestDeleteTitle(t *testing.T) {
	titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)
	titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

	creatorID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(creatorID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, creatorID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, creatorID, teamID); err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	cover, err := os.ReadFile("./test_data/title_cover.png")
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, creatorID, authorID, testhelpers.CreateTitleOptions{Cover: cover, Collection: titlesCovers})
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.TranslateTitle(env.DB, teamID, titleID); err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, env.NotificationsClient, titlesCovers, titlesOnModerationCovers)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.DELETE("/titles/:id", h.DeleteTitle)

	url := fmt.Sprintf("/titles/%d", titleID)
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

func TestEditTitle(t *testing.T) {
	titlesOnModerationCovers := env.MongoDB.Collection(mongodb.TitlesOnModerationCoversCollection)

	creatorID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(creatorID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, creatorID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, creatorID, teamID); err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, creatorID, authorID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.TranslateTitle(env.DB, teamID, titleID); err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, env.NotificationsClient, nil, titlesOnModerationCovers)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.POST("/titles/:id/edited", h.EditTitle)

	var authorName string
	env.DB.Raw("SELECT name FROM authors WHERE id = ?", authorID).Scan(&authorName) // Автор будет меняться на того же самого

	var genres []string
	h.DB.Raw("SELECT name FROM genres").Scan(&genres)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err = writer.WriteField("name", "newName"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("description", "newDescription"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("author", authorName); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(genres); i++ {
		if err = writer.WriteField("genres", genres[i]); err != nil {
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

func TestGetMostPopularTitles(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Genres: []string{"fighting"}})
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolume(env.DB, titleID, userID)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.ViewChapter(env.DB, userID, chapterID); err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.GET("/titles/most-popular", h.GetMostPopularTitles)

	req := httptest.NewRequest("GET", "/titles/most-popular?limit=10", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
}

func TestGetNewTitles(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Genres: []string{"fighting"}}); err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.GET("/titles/new", h.GetNewTitles)

	req := httptest.NewRequest("GET", "/titles/new?limit=5", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf(w.Body.String())
	}
}

func TestGetRecentlyUpdatedTitles(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Genres: []string{"fighting"}}); err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.GET("/titles/recently-updated", h.GetRecentlyUpdatedTitles)

	req := httptest.NewRequest("GET", "/titles/recently-updated?limit=10", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestGetTitleCover(t *testing.T) {
	titlesCovers := env.MongoDB.Collection(mongodb.TitlesCoversCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	cover, err := os.ReadFile("./test_data/title_cover.png")
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Cover: cover, Collection: titlesCovers})
	if err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, titlesCovers, nil)

	r := gin.New()
	r.GET("/titles/:id/cover", h.GetTitleCover)

	url := fmt.Sprintf("/titles/%d/cover", titleID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestGetTitle(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Genres: []string{"fighting"}})
	if err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.GET("/titles/:id", h.GetTitle)

	url := fmt.Sprintf("/titles/%d", titleID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestQuitTranslatingTitle(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
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

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{TeamID: teamID})
	if err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.PATCH("/titles/:id/quit-translating", h.QuitTranslatingTitle)

	url := fmt.Sprintf("/titles/%d/quit-translating", titleID)
	req := httptest.NewRequest("PATCH", url, nil)

	req.AddCookie(&http.Cookie{ // Надо будет это в хэлпер добавить
		Name:  "mangacage_token",
		Value: tokenString,
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestSubscribeToTitle(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID)
	if err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.POST("/titles/:id/subscriptions", h.SubscribeToTitle)

	url := fmt.Sprintf("/titles/%d/subscriptions", titleID)
	req := httptest.NewRequest("POST", url, nil)

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

func TestTranslateTitle(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
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

	authorID, err := testhelpers.CreateAuthor(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID)
	if err != nil {
		t.Fatal(err)
	}

	h := titles.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.PATCH("/titles/:id/translate", h.TranslateTitle)

	url := fmt.Sprintf("/titles/%d/translate", titleID)
	req := httptest.NewRequest("PATCH", url, nil)

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
