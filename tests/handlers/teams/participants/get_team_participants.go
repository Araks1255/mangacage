package participants

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTeamParticipantsScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":         GetTeamParticipantsSuccess(env),
		"invalid team id": GetTeamParticipantsWithInvalidTeamId(env),
		"wrong team id":   GetTeamParticipantsWithWrongTeamId(env),
		"no participants": GetTeamWithNoParticipantsParticipants(env),
	}
}

func GetTeamParticipantsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
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

		for i := 0; i < 5; i++ {
			if _, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Roles: []string{"typer"}, TeamID: teamID}); err != nil {
				t.Fatal(err)
			}
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

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) != 6 {
			t.Fatal("дошли не все участники")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id участника не дошел")
		}
		if _, ok := resp[0]["userName"]; !ok {
			t.Fatal("имя пользователя участника не дошло")
		}
		if roles, ok := resp[0]["roles"]; !ok || len(roles.([]any)) == 0 || resp[0]["roles"].([]any)[0].(string) != "typer" {
			t.Fatal("возникли проблемы с ролями")
		}
	}
}

func GetTeamParticipantsWithInvalidTeamId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		invalidTeamID := "||._.||"

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.GET("/teams/:id/participants", h.GetTeamParticipants)

		url := fmt.Sprintf("/teams/%s/participants", invalidTeamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTeamParticipantsWithWrongTeamId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		teamID := 9223372036854775807

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.GET("/teams/:id/participants", h.GetTeamParticipants)

		url := fmt.Sprintf("/teams/%d/participants", teamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTeamWithNoParticipantsParticipants(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := participants.NewHandler(env.DB)

		r := gin.New()
		r.GET("/teams/:id/participants", h.GetTeamParticipants)

		url := fmt.Sprintf("/teams/%d/participants", teamID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}
