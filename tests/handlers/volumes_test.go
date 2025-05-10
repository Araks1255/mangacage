package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/internal/testhelpers"
	"github.com/Araks1255/mangacage/pkg/handlers/volumes"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func TestCreateVolume(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleTranslatingByUserTeam(env.DB, userID, []string{"fighting"})
	if err != nil {
		t.Fatal(err)
	}

	h := volumes.NewHandler(env.DB, env.NotificationsClient)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.POST("/titles/:id/volumes", h.CreateVolume)

	body := gin.H{
		"name":        "volume",
		"description": "someDescription",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/titles/%d/volumes", titleID)
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

func TestDeleteVolume(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := volumes.NewHandler(env.DB, env.NotificationsClient)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.DELETE("/volumes/:id", h.DeleteVolume)

	url := fmt.Sprintf("/volumes/%d", volumeID)
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

func TestEditVolume(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := volumes.NewHandler(env.DB, env.NotificationsClient)

	r := gin.New()
	r.Use(middlewares.Auth(env.SecretKey))
	r.POST("/volumes/:id/edited", h.EditVolume)

	body := gin.H{
		"name":        "volume",
		"description": "someDescription",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/volumes/%d/edited", volumeID)
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

func TestGetTitleVolumes(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := testhelpers.CreateVolume(env.DB, titleID, userID); err != nil {
		t.Fatal(err)
	}

	h := volumes.NewHandler(env.DB, nil)

	r := gin.New()
	r.GET("/titles/:id/volumes", h.GetTitleVolumes)

	url := fmt.Sprintf("/titles/%d/volumes", titleID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestGetVolume(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	volumeID, err := testhelpers.CreateVolumeTranslatingByUserTeam(env.DB, userID)
	if err != nil {
		t.Fatal(err)
	}

	h := volumes.NewHandler(env.DB, env.NotificationsClient)

	r := gin.New()
	r.GET("/volumes/:id", h.GetVolume)

	url := fmt.Sprintf("/volumes/%d", volumeID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}
