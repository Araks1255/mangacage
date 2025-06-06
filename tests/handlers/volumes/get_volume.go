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

func GetGetVolumeScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":    GetVolumeSuccess(env),
		"wrong id":   GetVolumeWithWrongId(env),
		"invalid id": GetVolumeWithInvalidId(env),
	}
}

func GetVolumeSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
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

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if _, ok := resp["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp["createdAt"]; !ok {
			t.Fatal("время создания не дошло")
		}
		if _, ok := resp["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp["title"]; !ok {
			t.Fatal("тайтл не дошел")
		}
		if _, ok := resp["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
	}
}

func GetVolumeWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		volumeID := 9223372036854775807

		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/volumes/:id", h.GetVolume)

		url := fmt.Sprintf("/volumes/%d", volumeID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetVolumeWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := volumes.NewHandler(env.DB, env.NotificationsClient)

		r := gin.New()
		r.GET("/volumes/:id", h.GetVolume)

		req := httptest.NewRequest("GET", "/volumes/o_o", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
