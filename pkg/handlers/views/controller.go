package views

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine, secretKey string) {
	h := handler{DB: db}

	r.Static("/static", "static")
	r.LoadHTMLGlob("html/*.html")

	r.GET("/", h.ShowMainPage)
	r.GET("/titles/:id", h.ShowTitlePage)
	r.GET("/chapters/:id", middlewares.AuthOptional(secretKey), h.ShowChapterPage)
	r.GET("/teams/:id", h.ShowTeamPage)
	r.GET("/signup", h.ShowSignupPage)
	r.GET("/login", h.ShowLoginPage)
	r.GET("/titles", h.ShowTitlesCatalogPage)
	r.GET("/chapters", h.ShowChaptersCatalogPage)
	r.GET("/users/me", h.ShowMyProfilePage)
	r.GET("/users/me/moderation/titles", h.ShowTitlesOnModerationCatalog)
	r.GET("/users/me/moderation/titles/:id", h.ShowTitleOnModerationPage)
	r.GET("/users/me/moderation/chapters/:id", h.ShowChapterOnModerationPage)
	r.GET("/users/me/moderation/chapters", h.ShowChaptersOnModerationCatalog)
	r.GET("/users/me/moderation/teams", h.ShowTeamOnModerationCatalog)
	r.GET("/users/:id", h.ShowProfilePage)
	r.GET("/users/me/team-join-requests", h.ShowTeamJoinRequestsCatalog)
	r.GET("/users/me/titles-translate-requests", h.ShowTitleTranslateRequestsCatalog)
	r.GET("/teams", h.ShowTeamsCatalog)
	r.GET("/users", h.ShowUsersCatalog)
	r.GET("/privacy-policy", h.ShowPrivacyPolicy)
	r.GET("/terms-of-use", h.ShowTermsOfUse)
	r.GET("/copyright", h.ShowCopyright)
	r.GET("/rules", h.ShowRules)
	r.GET("/users/me/moderation/profile-changes", h.ShowProfileChangesCatalog)
}
