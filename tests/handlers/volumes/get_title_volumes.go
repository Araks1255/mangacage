package volumes

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/volumes"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTitleVolumesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":               GetTitleVolumesSuccess(env),
		"title without volumes": GetTitleWithoutVolumesVolumes(env),
		"wrong title id":        GetTitleVolumesWithWrongTitleId(env),
		"invalid title id":      GetTitleVolumesWithInvalidId(env),
		"invalid limit":         GetTitleVolumesWithInvalidLimit(env),
	}
}

func GetTitleVolumesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			if _, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/titles/:id/volumes", h.GetTitleVolumes)

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
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

		if len(resp) < 3 {
			t.Fatal("не все томы дошли")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("тайтд не дошёл")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошёл")
		}
	}
}

func GetTitleWithoutVolumesVolumes(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/titles/:id/volumes", h.GetTitleVolumes)

		url := fmt.Sprintf("/titles/%d/volumes", titleID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitleVolumesWithWrongTitleId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		wrongID := 9223372036854775807

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/titles/:id/volumes", h.GetTitleVolumes)

		url := fmt.Sprintf("/titles/%d/volumes", wrongID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitleVolumesWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/titles/:id/volumes", h.GetTitleVolumes)

		url := "/titles/Г_Г/volumes"
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitleVolumesWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/titles/:id/volumes", h.GetTitleVolumes)

		req := httptest.NewRequest("GET", "/titles/18/volumes?limit=O_O", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
