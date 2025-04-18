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
	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/pkg/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

// Teams
func TestCreateTeam(t *testing.T) {
	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	teamsOnModerationCovers := env.MongoDB.Collection(constants.TeamsOnModerationCoversCollection)
	teamsCovers := env.MongoDB.Collection(constants.TeamsCoversCollection)

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err = writer.WriteField("name", "teamName"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("description", "teamDescription"); err != nil {
		t.Fatal(err)
	}

	part, err := writer.CreateFormFile("cover", "file")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("./test_data/team_cover.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatal(err)
	}

	writer.Close()

	h := teams.NewHandler(env.DB, teamsOnModerationCovers, teamsCovers)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/teams", h.CreateTeam)

	req := httptest.NewRequest("POST", "/teams", &body)
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

func TestEditTeam(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID, Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamsOnModerationCovers := env.MongoDB.Collection(constants.TeamsOnModerationCoversCollection)
	teamsCovers := env.MongoDB.Collection(constants.TeamsCoversCollection)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err = writer.WriteField("name", "newName"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("description", "newDescription"); err != nil {
		t.Fatal(err)
	}

	part, err := writer.CreateFormFile("cover", "file")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("./test_data/team_cover.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatal(err)
	}

	writer.Close()

	h := teams.NewHandler(env.DB, teamsOnModerationCovers, teamsCovers)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))

	r.POST("/teams/:id", h.EditTeam)

	url := fmt.Sprintf("/teams/%d", teamID)
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

func TestGetTeamCover(t *testing.T) {
	teamsCovers := env.MongoDB.Collection(constants.TeamsCoversCollection)
	teamsOnModerationCovers := env.MongoDB.Collection(constants.TeamsOnModerationCoversCollection)

	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	h := teams.NewHandler(env.DB, teamsOnModerationCovers, teamsCovers)

	r := gin.New()
	r.GET("/teams/:id/cover", h.GetTeamCover)

	url := fmt.Sprintf("/teams/%d/cover", teamID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestGetTeam(t *testing.T) {
	teamsCovers := env.MongoDB.Collection(constants.TeamsCoversCollection)
	teamsOnModerationCovers := env.MongoDB.Collection(constants.TeamsOnModerationCoversCollection)

	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	h := teams.NewHandler(env.DB, teamsOnModerationCovers, teamsCovers)

	r := gin.New()
	r.GET("/teams/:id", h.GetTeam)

	url := fmt.Sprintf("/teams/%d", teamID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

// Join requests

func TestAcceptTeamJoinRequest(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	var teamLeaderID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&teamLeaderID)
	if teamLeaderID == 0 {
		t.Fatal("Тестовый юзер не найден")
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/teams/my/join-requests/:id/accept", h.AcceptTeamJoinRequest)

	url := fmt.Sprintf("/teams/my/join-requests/%d/accept", requestID)
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

func TestCancelTeamJoinRequest(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(candidateID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/teams/join-requests/:id", h.CancelTeamJoinRequest)

	url := fmt.Sprintf("/teams/join-requests/%d", requestID)
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

func TestDeclineTeamJoinRequest(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	var teamLeaderID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&teamLeaderID)
	if teamLeaderID == 0 {
		t.Fatal("Тестовый юзер не найден")
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/teams/my/join-requests/:id", h.DeclineTeamJoinRequest)

	url := fmt.Sprintf("/teams/my/join-requests/%d", requestID)
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

func TestGetMyTeamJoinRequests(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testhelpers.CreateTeamJoinRequest(env.DB, userID, teamID)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/teams/join-requests/my", h.GetMyTeamJoinRequests)

	req := httptest.NewRequest("GET", "/teams/join-requests/my", nil)

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

func TestGetTeamJoinRequestsOfMyTeam(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	var teamLeaderID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&teamLeaderID)
	if teamLeaderID == 0 {
		t.Fatal("Тестовый лидер команды не найден")
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.GET("/teams/my/join-requests", h.GetTeamJoinRequestsOfMyTeam)

	req := httptest.NewRequest("GET", "/teams/my/join-requests", nil)

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

func TestSubmitTeamJoinRequest(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	userID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

	body := map[string]string{
		"introductoryMessage": "message",
		"role":                "translater",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/teams/%d/join-requests", teamID)
	req := httptest.NewRequest("POST", url, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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

// Participants

func TestChangeParticipantRole(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	var teamLeaderID uint
	env.DB.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&teamLeaderID)
	if teamLeaderID == 0 {
		t.Fatal("Тестовый лидер команды не найден")
	}

	participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID, Roles: []string{"typer"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := participants.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.PATCH("/teams/my/participants/:id/role", h.ChangeParticipantRole)

	body := gin.H{
		"currentRole": "typer",
		"newRole":     "translater",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/teams/my/participants/%d/role", participantID)
	req := httptest.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

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

func TestLeaveTeam(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID, Roles: []string{"typer", "moder"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(participantID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	h := participants.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/teams/my/participants/me", h.LeaveTeam)

	req := httptest.NewRequest("DELETE", "/teams/my/participants/me", nil)

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

func TestGetTeamParticipants(t *testing.T) {
	var teamID uint
	env.DB.Raw("SELECT id FROM teams WHERE name = 'team_test'").Scan(&teamID)
	if teamID == 0 {
		t.Fatal("Тестовая команда не найдена")
	}

	h := participants.NewHandler(env.DB)

	r := gin.New()
	r.GET("/teams/:id/participants", h.GetTeamParticipants)

	url := fmt.Sprintf("/teams/%d/participants", teamID)
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}
