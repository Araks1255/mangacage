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

func GetGetVolumesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":      GetVolumesWithAllParamsSuccess(env),
		"success with query":      GetVolumesSuccessWithQuery(env),
		"success with pagination": GetVolumesWithPagination(env),
		"not found":               GetVolumesNotFound(env),
		"invalid order":           GetVolumesWithInvalidOrder(env),
	}
}

func GetVolumesWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
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

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateVolume(env.DB, titleID, teamID, userID); err != nil {
				t.Fatal(err)
			}

			if _, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := volumes.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/volumes", h.GetVolumes)

		url := fmt.Sprintf(
			"/volumes?sort=createdAt&order=desc&page=1&limit=20&titleId=%d&teamId=%d",
			titleID, teamID,
		)

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
			t.Fatal("неверное количество томов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["createdAt"]; !ok {
			t.Fatal("createdAt не дошел")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("название тайтла не дошло")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("название команды не дошло")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не дошел")
		}
	}
}

func GetVolumesSuccessWithQuery(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		volumeID, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var volumeName string
		if err := env.DB.Raw("SELECT name FROM volumes WHERE id = ?", volumeID).Scan(&volumeName).Error; err != nil {
			t.Fatal(err)
		}

		h := volumes.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/volumes", h.GetVolumes)

		url := fmt.Sprintf("/volumes?query=%s", volumeName)
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

		if len(resp) != 1 {
			t.Fatal("неверное количество томов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
	}
}

func GetVolumesWithPagination(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := volumes.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/volumes", h.GetVolumes)

		volumesIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/volumes?limit=1&page=%d&sort=createdAt", i)
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

			id, ok := resp[0]["id"].(float64)
			if !ok {
				t.Fatal("возникли проблемы с получением id")
			}

			volumesIDs[i-1] = uint(id)
		}

		if volumesIDs[0]-volumesIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией")
		}
	}
}

func GetVolumesNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := volumes.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/volumes", h.GetVolumes)

		req := httptest.NewRequest("GET", "/volumes?titleId=999999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetVolumesWithInvalidOrder(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateVolumeWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := volumes.NewHandler(env.DB, nil)

		r := gin.New()
		r.GET("/volumes", h.GetVolumes)

		req := httptest.NewRequest("GET", "/volumes?order=notvalid&sort=createdAt", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal(w.Body.String())
		}

		var resp []map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}

		if len(resp) < 2 {
			t.Fatal("возникли проблемы с количеством томов")
		}

		if uint(resp[0]["id"].(float64))-uint(resp[1]["id"].(float64)) != 1 { // При невалидном order должен выставиться desc
			t.Fatal("возникли проблемы с порядком томов")
		}
	}
}
