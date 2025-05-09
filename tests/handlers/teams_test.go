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
	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
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

	teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsOnModerationCoversCollection)

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

	h := teams.NewHandler(env.DB, teamsOnModerationCovers, nil)

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

	tokenString, err := testhelpers.GenerateTokenString(userID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsOnModerationCoversCollection)
	teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

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
	teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)
	teamsOnModerationCovers := env.MongoDB.Collection(mongodb.TeamsOnModerationCoversCollection)

	creatorID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	cover, err := os.ReadFile("./test_data/team_cover.png")
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, creatorID, testhelpers.CreateTeamOptions{Cover: cover, Collection: teamsCovers})
	if err != nil {
		t.Fatal(err)
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
	creatorID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
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

	h := teams.NewHandler(env.DB, nil, nil)

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
	teamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, teamLeaderID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, teamLeaderID, teamID); err != nil {
		t.Fatal(err)
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
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

	if w.Code != 200 {
		t.Fatal(w.Body.String())
	}
}

func TestCancelTeamJoinRequest(t *testing.T) {
	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(candidateID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, candidateID)
	if err != nil {
		t.Fatal(err)
	}

	requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
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
	teamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, teamLeaderID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, teamLeaderID, teamID); err != nil {
		t.Fatal(err)
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	requestID, err := testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
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
	userID, err := testhelpers.CreateUser(env.DB)
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

	_, err = testhelpers.CreateTeamJoinRequest(env.DB, userID, teamID)
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
	teamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, teamLeaderID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, teamLeaderID, teamID); err != nil {
		t.Fatal(err)
	}

	candidateID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testhelpers.CreateTeamJoinRequest(env.DB, candidateID, teamID)
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
	userID, err := testhelpers.CreateUser(env.DB)
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

	var roleID uint
	env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID)
	if roleID == 0 {
		t.Fatal("роль не найдена")
	}

	h := joinrequests.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/teams/:id/join-requests", h.SubmitTeamJoinRequest)

	body := gin.H{
		"introductoryMessage": "message",
		"roleId":              roleID,
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

	if w.Code != 201 {
		t.Fatal(w.Body.String())
	}
}

// Participants

func TestAddRoleToParticipant(t *testing.T) {
	participantID, err := testhelpers.CreateUser(env.DB)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, participantID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, participantID, teamID); err != nil {
		t.Fatal(err)
	}

	teamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID, Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var roleID uint
	env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID) // Эта роль, как и остальные, необходима для корректной работы бэкенда. Так что тут можно просто достать из бд
	if roleID == 0 {
		t.Fatal("Роль не найдена")
	}

	h := participants.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.POST("/teams/my/participants/:id/roles", h.AddRoleToParticipant)

	body := map[string]uint{
		"roleId": roleID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
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

func TestDeleteParticipantRole(t *testing.T) {
	participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}})
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, participantID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, participantID, teamID); err != nil {
		t.Fatal(err)
	}

	teamLeaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{TeamID: teamID, Roles: []string{"team_leader"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(teamLeaderID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	var roleID uint
	env.DB.Raw("SELECT id FROM roles WHERE name = 'typer'").Scan(&roleID) // Эта роль, как и остальные, необходима для корректной работы бэкенда. Так что тут можно просто достать из бд
	if roleID == 0 {
		t.Fatal("Роль не найдена")
	}

	h := participants.NewHandler(env.DB)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(env.SecretKey))
	r.DELETE("/teams/my/participants/:id/roles", h.DeleteParticipantRole)

	body := map[string]uint{
		"roleId": roleID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("/teams/my/participants/%d/roles", participantID)
	req := httptest.NewRequest("DELETE", url, bytes.NewBuffer(jsonBody))
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
	participantID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer", "moder"}})
	if err != nil {
		t.Fatal(err)
	}

	tokenString, err := testhelpers.GenerateTokenString(participantID, env.SecretKey)
	if err != nil {
		t.Fatal(err)
	}

	teamID, err := testhelpers.CreateTeam(env.DB, participantID)
	if err != nil {
		t.Fatal(err)
	}

	if err = testhelpers.AddUserToTeam(env.DB, participantID, teamID); err != nil {
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
	userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}})
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
