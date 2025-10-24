package chapters

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetGetChaptersScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success all params":          GetChaptersWithAllParamsSuccess(env),
		"success with query":          GetChaptersWithQuerySuccess(env),
		"success with pagination":     GetChaptersWithPaginationSuccess(env),
		"success my favorites":        GetMyFavoriteChaptersSuccess(env),
		"success favorited by user":   GetChaptersFavoritedByUserSuccess(env),
		"favorites unauthorized":      GetMyFavoriteChaptersUnauthorized(env),
		"favorited by invisible user": GetChaptersFavoritedByInvisibleUser(env),
		"not found":                   GetChaptersNotFound(env),
		"invalid order":               GetChaptersWithInvalidOrder(env),
	}
}

func GetChaptersWithAllParamsSuccess(env testenv.Env) func(*testing.T) {
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
			_, err := testhelpers.CreateChapter(
				env.DB, titleID, teamID, userID,
				testhelpers.CreateChapterOptions{
					Pages: make([][]byte, 5+i),
					Views: uint(5 + i),
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		url := fmt.Sprintf(
			"/chapters?sort=createdAt&order=desc&page=1&limit=20&volume=0&titleId=%d&teamId=%d&numberOfPagesFrom=5&numberOfPagesTo=1000&viewsFrom=5&viewsTo=1000",
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
			t.Fatal("неверное количество глав")
		}

		if _, ok := resp[0]["id"]; !ok {
			t.Fatal("id не дошел")
		}
		if _, ok := resp[0]["name"]; !ok {
			t.Fatal("название не дошло")
		}
		if _, ok := resp[0]["team"]; !ok {
			t.Fatal("название команды не дошло")
		}
		if _, ok := resp[0]["teamId"]; !ok {
			t.Fatal("id команды не дошел")
		}
		if _, ok := resp[0]["titleId"]; !ok {
			t.Fatal("id тайтла не дошел")
		}
		if _, ok := resp[0]["title"]; !ok {
			t.Fatal("название тайтла не дошло")
		}
	}
}

func GetChaptersWithQuerySuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		var chapterName string
		if err := env.DB.Raw("SELECT name FROM chapters WHERE id = ?", chapterID).Scan(&chapterName).Error; err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		url := fmt.Sprintf("/chapters?query=%s", chapterName)
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
			t.Fatal("неверное количество глав")
		}
	}
}

func GetChaptersWithPaginationSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		chaptersIDs := make([]uint, 2)

		for i := 1; i <= 2; i++ {
			url := fmt.Sprintf("/chapters?limit=1&page=%d&sort=createdAt", i)
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

			chaptersIDs[i-1] = uint(id)
		}

		if chaptersIDs[0]-chaptersIDs[1] != 1 {
			t.Fatal("возникли проблемы с пагинацией")
		}
	}
}

func GetMyFavoriteChaptersSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		req := httptest.NewRequest("GET", "/chapters?myFavorites=true", nil)

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
			t.Fatal("пришло не столько глав, сколько ожидалось")
		}
		if uint(resp[0]["id"].(float64)) != chapterID {
			t.Fatal("пришла не та глава")
		}
	}
}

func GetChaptersFavoritedByUserSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB, testhelpers.CreateUserOptions{Visible: true})
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		url := fmt.Sprintf("/chapters?favoritedBy=%d", userID)
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
			t.Fatal("пришло не столько глав, сколько ожидалось")
		}
		if uint(resp[0]["id"].(float64)) != chapterID {
			t.Fatal("пришла не та глава")
		}
	}
}

func GetMyFavoriteChaptersUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		req := httptest.NewRequest("GET", "/chapters?myFavorites=true", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChaptersFavoritedByInvisibleUser(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		chapterID, err := testhelpers.CreateChapterWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		if err := testhelpers.AddChapterToFavorites(env.DB, userID, chapterID); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		url := fmt.Sprintf("/chapters?favoritedBy=%d", userID)
		req := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChaptersNotFound(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
			t.Fatal(err)
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		req := httptest.NewRequest("GET", "/chapters?numberOfPagesFrom=1000", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func GetChaptersWithInvalidOrder(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			if _, err := testhelpers.CreateChapterWithDependencies(env.DB, userID); err != nil {
				t.Fatal(err)
			}
		}

		h := chapters.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.AuthOptional(env.SecretKey))
		r.GET("/chapters", h.GetChapters)

		req := httptest.NewRequest("GET", "/chapters?order=0_0&sort=createdAt", nil)

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
			t.Fatal("возникли проблемы с количеством глав")
		}

		if uint(resp[0]["id"].(float64))-uint(resp[1]["id"].(float64)) != 1 { // При невалидном order должен выставиться desc
			t.Fatal("возникли проблемы с порядком глав")
		}
	}
}
