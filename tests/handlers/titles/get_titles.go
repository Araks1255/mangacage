package titles

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetTitlesScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":         GetTitlesWithAllParamsSuccess(env),
		"success with query":         GetTitlesWithQuerySuccess(env),
		"success with pagination":    GetTitlesWithPaginationSuccess(env),
		"success my favorites":       GetMyFavoriteTitlesSuccess(env),
		"favorited by user success":  GetTitlesFavoritedByUserSuccess(env),
		"favorite by invisible user": GetTitlesFavoritedByInvisibleUser(env),
		"my favorites unauthorized":  GetMyFavoriteTitlesUnauthorized(env),
		"not found":                  GetTitlesNotFound(env),
		"invalid order":              GetTitlesWithInvalidOrder(env),
	}
}

func GetTitlesWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		authorID, err := testhelpers.CreateAuthor(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		teamID, err := testhelpers.CreateTeam(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			_, err := testhelpers.CreateTitle(
				env.DB, userID, authorID, testhelpers.CreateTitleOptions{
					TeamID: teamID, Views: uint(5 + i),
					Genres: []string{"action", "fighting"},
					Tags:   []string{"maids", "japan"},
				})

			if err != nil {
				t.Fatal(err)
			}

			if _, err = testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		url := fmt.Sprintf(
			"/titles?sort=createdAt&order=desc&page=1&limit=20&publishingStatus=ongoing&translatingStatus=ongoing&type=manga&authorId=%d&teamId=%d&yearFrom=1000&yearTo=3000&viewsFrom=5&viewsTo=1000&genres=action&genres=fighting&tags=maids&tags=japan",
			authorID, teamID,
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
			t.Fatal("неверное количество тайтлов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
		if _, ok := resp[0]["views"]; !ok {
			t.Fatal("возникли проблемы с просмотрами")
		}
	}
}

func GetTitlesWithQuerySuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var titleName string
		if err := env.DB.Raw("SELECT name FROM titles WHERE id = ?", titleID).Scan(&titleName).Error; err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		url := fmt.Sprintf("/titles?query=%s", titleName)
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
			t.Fatal("неверное количество тайтлов")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["authorId"]; !ok {
			t.Fatal("id автора не дошел")
		}
	}
}

func GetTitlesWithPaginationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		titlesIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/titles?limit=1&page=%d&sort=createdAt", i)
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

			titlesIDs[i-1] = uint(id)
		}

		if titlesIDs[0]-titlesIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией", titlesIDs)
		}
	}
}

func GetMyFavoriteTitlesSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		req := httptest.NewRequest("GET", "/titles?myFavorites=true", nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

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
			t.Fatal("пришло не столько тайтлов, сколько ожидалось")
		}
		if uint(resp[0]["id"].(float64)) != titleID {
			t.Fatal("пришел не тот тайтл")
		}
	}
}
func GetTitlesFavoritedByUserSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Visible: true})
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		url := fmt.Sprintf("/titles?favoritedBy=%d", userID)
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
			t.Fatal("пришло не столько тайтлов, сколько ожидалось")
		}
		if uint(resp[0]["id"].(float64)) != titleID {
			t.Fatal("пришел не тот тайтл")
		}
	}
}

func GetTitlesFavoritedByInvisibleUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddTitleToFavorites(env.DB, userID, titleID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		url := fmt.Sprintf("/titles?favoritedBy=%d", userID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetMyFavoriteTitlesUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		req := httptest.NewRequest("GET", "/titles?myFavorites=true", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitlesNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		req := httptest.NewRequest("GET", "/titles?yearTo=1", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetTitlesWithInvalidOrder(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateTitleWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/titles", h.GetTitles)

		req := httptest.NewRequest("GET", "/titles?order=4554&sort=createdAt", nil)

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
			t.Fatal("возникли проблемы с количеством тайтлов")
		}

		if uint(resp[0]["id"].(float64))-uint(resp[1]["id"].(float64)) != 1 { // При невалидном order должен выставиться desc
			t.Fatal("возникли проблемы с порядком тайтлов")
		}
	}
}
