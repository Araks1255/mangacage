package titles

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetSubscribeToTitleScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"success":      SubscribeToTitleSuccess(env),
		"unauthorized": SubscribeToTitleUnauthorized(env),
		"twice":        SubscribeToTitleTwice(env),
		"wrong id":     SubscribeToTitleWithWrongId(env),
		"invalid id":   SubscribeToTitleWithInvalidId(env),
	}
}

func SubscribeToTitleSuccess(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles/:id/subscriptions", h.SubscribeToTitle)

		url := fmt.Sprintf("/titles/%d/subscriptions", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}
	}
}

func SubscribeToTitleUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles/:id/subscriptions", h.SubscribeToTitle)

		req := httptest.NewRequest("POST", "/titles/18/subscriptions", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}
func SubscribeToTitleTwice(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID, err := testhelpers.CreateTitleWithDependencies(env.DB, userID)
		if err != nil {
			t.Fatal(err)
		}

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles/:id/subscriptions", h.SubscribeToTitle)

		url := fmt.Sprintf("/titles/%d/subscriptions", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Fatal(w.Body.String())
		}

		req2 := httptest.NewRequest("POST", url, nil)
		req2.AddCookie(cookie)

		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		if w2.Code != 409 {
			t.Fatal(w2.Body.String())
		}
	}
}

func SubscribeToTitleWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		titleID := 9223372036854775807

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles/:id/subscriptions", h.SubscribeToTitle)

		url := fmt.Sprintf("/titles/%d/subscriptions", titleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 404 {
			t.Fatal(w.Body.String())
		}
	}
}

func SubscribeToTitleWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		invalidTitleID := "@_@"

		h := titles.NewHandler(env.DB, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.POST("/titles/:id/subscriptions", h.SubscribeToTitle)

		url := fmt.Sprintf("/titles/%s/subscriptions", invalidTitleID)
		req := httptest.NewRequest("POST", url, nil)

		cookie, err := testhelpers.CreateCookieWithToken(userID, env.SecretKey)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(cookie)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 400 {
			t.Fatal(w.Body.String())
		}
	}
}
