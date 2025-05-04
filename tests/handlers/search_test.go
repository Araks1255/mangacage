package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/internal/testhelpers"
	"github.com/Araks1255/mangacage/pkg/handlers/search"
	"github.com/gin-gonic/gin"
)

func TestSearchAuthors(t *testing.T) {
	authorID, err := testhelpers.CreateAuthor(env.DB, testhelpers.CreateAuthorOptions{Genres: []string{"fighting"}})
	if err != nil {
		t.Fatal(err)
	}
	var authorName string
	env.DB.Raw("SELECT name FROM authors WHERE id = ?", authorID).Scan(&authorName)

	h := search.NewHandler(env.DB)

	r := gin.New()
	r.GET("/search", h.Search)

	query := authorName[:5]
	url := fmt.Sprintf("/search?type=authors&query=%s&limit=10", query)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestSearchChapters(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	chapterID, err := testhelpers.CreateChapter(env.DB, volumeID, userID)
	if err != nil {
		t.Fatal(err)
	}
	var chapterName string
	env.DB.Raw("SELECT name FROM chapters WHERE id = ?", chapterID).Scan(&chapterName)

	h := search.NewHandler(env.DB)
	r := gin.New()
	r.GET("/search", h.Search)

	query := chapterName[:5]
	url := fmt.Sprintf("/search?type=chapters&query=%s&limit=10", query)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestSearchTeams(t *testing.T) {
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

	var teamName string
	env.DB.Raw("SELECT name FROM teams WHERE id = ?", teamID).Scan(&teamName)

	h := search.NewHandler(env.DB)
	r := gin.New()
	r.GET("/search", h.Search)

	query := teamName[:5]
	url := fmt.Sprintf("/search?type=teams&query=%s&limit=10", query)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestSearchTitles(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	authorID, err := testhelpers.CreateAuthor(env.DB)

	titleID, err := testhelpers.CreateTitle(env.DB, userID, authorID)
	if err != nil {
		t.Fatal(err)
	}
	var titleName string
	env.DB.Raw("SELECT name FROM titles WHERE id = ?", titleID).Scan(&titleName)

	h := search.NewHandler(env.DB)
	r := gin.New()
	r.GET("/search", h.Search)

	query := titleName[:5]
	url := fmt.Sprintf("/search?type=titles&query=%s&limit=10", query)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestSearchVolumes(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}
	var volumeName string
	env.DB.Raw("SELECT name FROM volumes WHERE id = ?", volumeID).Scan(&volumeName)

	h := search.NewHandler(env.DB)
	r := gin.New()
	r.GET("/search", h.Search)

	query := volumeName[:5]
	url := fmt.Sprintf("/search?type=volumes&query=%s&limit=10", query)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}
