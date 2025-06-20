package teams

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/pkg/constants/mongodb"
	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTeamCoverScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         GetTeamCoverSuccess(env),
		"wrong team id":   GetTeamCoverWithWrongTeamId(env),
		"invalid team id": GetTeamCoverWithInvalidTeamId(env),
	}
}

func GetTeamCoverSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		cover, err := os.ReadFile("./test_data/team_cover.png")
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID, testhelpers.CreateTeamOptions{Cover: cover, Collection: teamsCovers})
		if err != nil {
			t.Fatal(err)
		}

		h := teams.NewHandler(env.DB, teamsCovers)

		r := gin.New()
		r.GET("/teams/:id/cover", h.GetTeamCover)

		url := fmt.Sprintf("/teams/%d/cover", teamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		if len(w.Body.Bytes()) != len(cover) {
			t.Fatal("обложка не отправилась")
		}
	}
}

func GetTeamCoverWithWrongTeamId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		teamID := 9223372036854775807

		h := teams.NewHandler(env.DB, teamsCovers)

		r := gin.New()
		r.GET("/teams/:id/cover", h.GetTeamCover)

		url := fmt.Sprintf("/teams/%d/cover", teamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTeamCoverWithInvalidTeamId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamsCovers := env.MongoDB.Collection(mongodb.TeamsCoversCollection)

		invalidTeamID := "8)"

		h := teams.NewHandler(env.DB, teamsCovers)

		r := gin.New()
		r.GET("/teams/:id/cover", h.GetTeamCover)

		url := fmt.Sprintf("/teams/%s/cover", invalidTeamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
