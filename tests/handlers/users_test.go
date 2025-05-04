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
	moderationHelpers "github.com/Araks1255/mangacage/internal/testhelpers/moderation"
	"github.com/Araks1255/mangacage/pkg/constants"
	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/handlers/users/favorites"
	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

// Users
func TestEditProfile(t *testing.T) {
	usersOnModerationProfilePictures := env.MongoDB.Collection(constants.UsersOnModerationProfilePicturesCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := users.NewHandler(env.DB, env.NotificationsClient, nil, usersOnModerationProfilePictures)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/users/me/edited", h.EditProfile)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("userName", "newUserName"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("aboutYourself", "newAbout"); err != nil {
		t.Fatal(err)
	}

	part, err := writer.CreateFormFile("profilePicture", "file")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("./test_data/user_profile_picture.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", "/users/me/edited", &body)
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

func TestGetMyProfilePicture(t *testing.T) {
	usersProfilePictures := env.MongoDB.Collection(constants.UsersProfilePicturesCollection)

	data, err := os.ReadFile("./test_data/user_profile_picture.png")
	if err != nil {
		t.Fatal(err)
	}

	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{ProfilePicture: data, Collection: usersProfilePictures})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := users.NewHandler(env.DB, env.NotificationsClient, usersProfilePictures, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/profile-picture", h.GetMyProfilePicture)

	req := httptest.NewRequest("GET", "/users/me/profile-picture", nil)

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

func TestGetMyProfile(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := users.NewHandler(env.DB, env.NotificationsClient, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me", h.GetMyProfile)

	req := httptest.NewRequest("GET", "/users/me", nil)

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

// Favorites

func TestAddTitleToFavorites(t *testing.T) {
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

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/favorites/titles/:id", h.AddTitleToFavorites)

	url := fmt.Sprintf("/favorites/titles/%d", titleID)
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

func TestGetFavoriteTitles(t *testing.T) {
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

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID, testhelpers.CreateTitleOptions{Genres: []string{"fighting"}})
	if err != nil {
		t.Fatal(err)
	}

	if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/favorites/titles", h.GetFavoriteTitles)

	req := httptest.NewRequest("GET", "/favorites/titles", nil)

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

func TestDeleteTitleFromFavorites(t *testing.T) {
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

	if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/favorites/titles/:id", h.DeleteTitleFromFavorites)

	url := fmt.Sprintf("/favorites/titles/%d", titleID)
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

func TestAddGenreToFavorites(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var genreID uint
	env.DB.Raw("SELECT id FROM genres WHERE name = 'fighting'").Scan(&genreID)

	h := favorites.NewHandler(env.DB)
	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/favorites/genres/:id", h.AddGenreToFavorites)

	url := fmt.Sprintf("/favorites/genres/%d", genreID)
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

func TestGetFavoriteGenres(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var genreID uint
	env.DB.Raw("SELECT id FROM genres WHERE name = 'fighting'").Scan(&genreID)

	if err := testhelpers.AddGenreToFavorites(env.DB, userID, genreID); err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/favorites/genres", h.GetFavoriteGenres)

	req := httptest.NewRequest("GET", "/favorites/genres", nil)

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

func TestDeleteGenreFromFavorites(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var genreID uint
	env.DB.Raw("SELECT id FROM genres WHERE name = 'fighting'").Scan(&genreID)

	if err := testhelpers.AddGenreToFavorites(env.DB, userID, genreID); err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/favorites/genres/:id", h.DeleteGenreFromFavorites)

	url := fmt.Sprintf("/favorites/genres/%d", genreID)
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

func TestAddChapterToFavorites(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/favorites/chapters/:id", h.AddChapterToFavorites)

	url := fmt.Sprintf("/favorites/chapters/%d", chapterID)
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

func TestGetFavoriteChapters(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/favorites/chapters", h.GetFavoriteChapters)

	req := httptest.NewRequest("GET", "/favorites/chapters", nil)

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

func TestDeleteChapterFromFavorites(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
		t.Fatal(err)
	}

	h := favorites.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/favorites/chapters/:id", h.DeleteChapterFromFavorites)

	url := fmt.Sprintf("/favorites/chapters/%d", chapterID)
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

// Moderation

func TestCancelAppealForChapterOnModeration(t *testing.T) {
	chaptersOnModerationsPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	chapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, chaptersOnModerationsPages, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/users/me/moderation/chapters/:id", h.CancelAppealForChapterModeration)

	url := fmt.Sprintf("/users/me/moderation/chapters/%d", chapterOnModerationID)
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

func TestGetMyChapterOnModerationPage(t *testing.T) {
	chaptersOnModerationPages := env.MongoDB.Collection(constants.ChaptersOnModerationPagesCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	pages := make([][]byte, 1, 1)

	pages[0], err = os.ReadFile("./test_data/chapter_page.png")
	if err != nil {
		t.Fatal(err)
	}

	chapterOnModerationID, err := moderationHelpers.CreateChapterOnModerationWithDependencies(
		env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Pages: pages, Collection: chaptersOnModerationPages},
	)

	if err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, chaptersOnModerationPages, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/chapters/:id/page/:page", h.GetMyChapterOnModerationPage)

	url := fmt.Sprintf("/users/me/moderation/chapters/%d/page/%d", chapterOnModerationID, 0)
	req := httptest.NewRequest("GET", url, nil)

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

func TestGetMyEditedChaptersOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateChapterOnModerationWithDependencies(env.DB, userID, moderationHelpers.CreateChapterOnModerationWithDependenciesOptions{Edited: true}); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/chapters/edited", h.GetMyEditedChaptersOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/chapters/edited", nil)

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

func TestGetMyNewChaptersOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateChapterOnModerationWithDependencies(env.DB, userID); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/chapters/new", h.GetMyNewChaptersOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/chapters/new", nil)

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

func TestCancelAppealForProfileChanges(t *testing.T) {
	usersOnModerationProfilePictures := env.MongoDB.Collection(constants.UsersOnModerationProfilePicturesCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateUserOnModeration(env.DB, moderationHelpers.CreateUserOnModerationOptions{ExistingID: userID}); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, usersOnModerationProfilePictures)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/users/me/moderation/profile/edited", h.CancelAppealForProfileChanges)

	req := httptest.NewRequest("DELETE", "/users/me/moderation/profile/edited", nil)

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

func TestGetMyProfilePictureOnModeration(t *testing.T) {
	usersOnModerationProfilePictures := env.MongoDB.Collection(constants.UsersOnModerationProfilePicturesCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("./test_data/user_profile_picture.png")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateUserOnModeration(
		env.DB, moderationHelpers.CreateUserOnModerationOptions{ExistingID: userID, ProfilePicture: data, Collection: usersOnModerationProfilePictures},
	); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, usersOnModerationProfilePictures)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/profile/picture", h.GetMyProfilePictureOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/profile/picture", nil)

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

func TestGetMyProfileChangesOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateUserOnModeration(env.DB, moderationHelpers.CreateUserOnModerationOptions{ExistingID: userID}); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/profile/edited", h.GetMyProfileChangesOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/profile/edited", nil)

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

func TestCancelAppealForTitleModeration(t *testing.T) {
	titlesOnModerationCovers := env.MongoDB.Collection(constants.TitlesOnModerationCoversCollection)

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, titlesOnModerationCovers, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/users/me/moderation/titles/:id", h.CancelAppealForTitleModeration)

	url := fmt.Sprintf("/users/me/moderation/titles/%d", titleOnModerationID)
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

func TestGetMyTitleOnModerationCover(t *testing.T) {
	titlesOnModerationCovers := env.MongoDB.Collection(constants.TitlesOnModerationCoversCollection)

	data, err := os.ReadFile("./test_data/title_cover.png")
	if err != nil {
		t.Fatal(err)
	}

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleOnModerationID, err := moderationHelpers.CreateTitleOnModeration(env.DB, userID, moderationHelpers.CreateTitleOnModerationOptions{Cover: data, Collection: titlesOnModerationCovers})
	if err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, titlesOnModerationCovers, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/titles/:id/cover", h.GetMyTitleOnModerationCover)

	url := fmt.Sprintf("/users/me/moderation/titles/%d/cover", titleOnModerationID)
	req := httptest.NewRequest("GET", url, nil)

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

func TestGetMyNewTitlesOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateTitleOnModeration(env.DB, userID); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/titles/new", h.GetMyNewTitlesOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/titles/new", nil)

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

func TestGetMyEditedTitlesOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateTitleOnModeration(env.DB, userID, moderationHelpers.CreateTitleOnModerationOptions{ExistingID: titleID}); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/titles/edited", h.GetMyEditedTitlesOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/titles/edited", nil)

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

func TestCancelAppealForVolumeOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	volumeOnModerationID, err := moderationHelpers.CreateVolumeOnModeration(env.DB, titleID, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/users/me/moderation/volumes/:id", h.CancelAppealForVolumeModeration)

	url := fmt.Sprintf("/users/me/moderation/volumes/%d", volumeOnModerationID)
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

func TestGetMyNewVolumesOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateVolumeOnModeration(env.DB, titleID, userID); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/volumes/new", h.GetMyNewVolumesOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/volumes/new", nil)

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

func TestGetMyEditedVolumesOnModeration(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolume(env.DB, titleID, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := moderationHelpers.CreateVolumeOnModeration(env.DB, titleID, userID, moderationHelpers.CreateVolumeOnModerationOptions{ExistingID: volumeID}); err != nil {
		t.Fatal(err)
	}

	h := moderation.NewHandler(env.DB, nil, nil, nil)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/users/me/moderation/volumes/edited", h.GetMyEditedVolumesOnModeration)

	req := httptest.NewRequest("GET", "/users/me/moderation/volumes/edited", nil)

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
