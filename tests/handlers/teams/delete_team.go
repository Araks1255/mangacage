package teams

import (
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/testhelpers/moderation"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetDeleteTeamScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         DeleteTeamSuccess(env),
		"unauthorized":    DeleteTeamByUnauthorizedUser(env),
		"non team leader": DeleteTeamByNonTeamLeader(env),
	}
}

func DeleteTeamSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		leaderID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"team_leader"}})
		if err != nil {
			t.Fatal(err)
		}

		cover := make([]byte, 1<<20)
		teamID, err := testhelpers.CreateTeam(env.DB, leaderID, testhelpers.CreateTeamOptions{Cover: cover, Collection: teamsCovers})
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, leaderID, teamID); err != nil {
			t.Fatal(err)
		}

		participantID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddUserToTeam(env.DB, participantID, teamID); err != nil {
			t.Fatal(err)
		}

		if _, err := moderation.CreateTeamOnModeration(
			env.DB, leaderID, moderation.CreateTeamOnModerationOptions{ExistingID: teamID, Cover: cover, Collection: teamsCovers},
		); err != nil {
			t.Fatal(err)
		}

		h := teams.NewHandler(env.DB, teamsCovers, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my", h.DeleteTeam)

		req := httptest.NewRequest("DELETE", "/teams/my", nil)

		cookie, err := testhelpers.CreateCookieWithToken(leaderID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteTeamByUnauthorizedUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := teams.NewHandler(env.DB, nil, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my", h.DeleteTeam)

		req := httptest.NewRequest("DELETE", "/teams/my", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func DeleteTeamByNonTeamLeader(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := teams.NewHandler(env.DB, nil, env.NotificationsClient)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey), middlewares.RequireRoles(env.DB, []string{"team_leader"}))
		r.DELETE("/teams/my", h.DeleteTeam)

		req := httptest.NewRequest("DELETE", "/teams/my", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 403 {
			t.Fatal(w.Body.String())
		}
	}
}
