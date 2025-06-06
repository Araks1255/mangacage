package moderation

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/Araks1255/mangacage/testhelpers"
	"github.com/Araks1255/mangacage/tests/testenv"
	"github.com/gin-gonic/gin"
)

func GetCancelAppealForModerationScenarios(env testenv.Env) map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		"titles":       CancelAppealForTitleModeration(env),
		"volumes":      CancelAppealForVolumeModeration(env),
		"chapters":     CancelAppealForChapterModeration(env),
		"teams":        CancelAppealForTeamModeration(env),
		"wrong entity": CancelAppealForModerationWithWrongEntity(env),
		"invalid id":   CancelAppealForModerationWithInvalidId(env),
		"wrong id":     CancelAppealForModerationWithWrongId(env),
		"unauthorized": CancelAppealForModerationUnauthorized(env),
	}
}

func CancelAppealForTitleModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		scenarios := GetCancelAppealForTitleModerationScenarios(env)

		for name, scenario := range scenarios {
			t.Run(name, scenario)
		}
	}
}

func CancelAppealForVolumeModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		scenarios := GetCancelAppealForVolumeModerationScenarios(env)

		for name, scenario := range scenarios {
			t.Run(name, scenario)
		}
	}
}

func CancelAppealForChapterModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		scenarios := GetCancelAppealForChapterModerationScenarios(env)

		for name, scenario := range scenarios {
			t.Run(name, scenario)
		}
	}
}

func CancelAppealForTeamModeration(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		scenarios := GetCancelAppealForTeamModerationScenarios(env)

		for name, scenario := range scenarios {
			t.Run(name, scenario)
		}
	}
}

func CancelAppealForModerationWithWrongEntity(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		req := httptest.NewRequest("DELETE", "/users/me/moderation/joinrequests/18", nil)

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

func CancelAppealForModerationWithInvalidId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		req := httptest.NewRequest("DELETE", "/users/me/moderation/titles/`-`", nil)

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

func CancelAppealForModerationWithWrongId(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		userID, err := testhelpers.CreateUser(env.DB)
		if err != nil {
			t.Fatal(err)
		}

		entityID := 9223372036854775807

		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		url := fmt.Sprintf("/users/me/moderation/titles/%d", entityID)
		req := httptest.NewRequest("DELETE", url, nil)

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

func CancelAppealForModerationUnauthorized(env testenv.Env) func(*testing.T) {
	return func(t *testing.T) {
		h := moderation.NewHandler(env.DB, nil, nil, nil, nil)

		r := gin.New()
		r.Use(middlewares.Auth(env.SecretKey))
		r.DELETE("/users/me/moderation/:entity/:id", h.CancelAppealForModeration)

		req := httptest.NewRequest("DELETE", "/users/me/moderation/titles/18", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 401 {
			t.Fatal(w.Body.String())
		}
	}
}
