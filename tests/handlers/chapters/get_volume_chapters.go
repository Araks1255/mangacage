package chapters

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetVolumeChaptersScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":           GetVolumeChaptersSuccess(env),
		"invalid volume id": GetVolumeChaptersWithInvalidVolumeId(env),
		"invalid limit":     GetVolumeChaptersWithInvalidLimit(env),
		"wrong volume id":   GetVolumeChaptersWithWrongVolumeId(env),
	}
}

func GetVolumeChaptersSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateChapter(env.DB, volumeID, teamID, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/volume/:id/chapters", h.GetVolumeChapters)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
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

		if len(resp) != 2 {
			t.Fatal("отобразились не все главы")
		}
		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("не отправился id")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("не отправилось название")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("не отправилось время создания")
		}
	}
}

func GetVolumeChaptersWithInvalidVolumeId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		volumeID := "`~`"
		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/volume/:id/chapters", h.GetVolumeChapters)

		url := fmt.Sprintf("/volume/%s/chapters", volumeID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetVolumeChaptersWithInvalidLimit(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/volume/:id/chapters", h.GetVolumeChapters)

		req := httptest.NewRequest("GET", "/volume/18/chapters?limit=U_U", nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetVolumeChaptersWithWrongVolumeId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		volumeID := 9223372036854775807

		h := chapters.NewHandler(env.DB, nil, nil, nil)

		r := gin.New()
		r.GET("/volume/:id/chapters", h.GetVolumeChapters)

		url := fmt.Sprintf("/volume/%d/chapters", volumeID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}
