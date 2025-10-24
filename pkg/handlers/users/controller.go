package users

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type handler struct {
	DB                  *gorm.DB
	PathToMediaDir      string
	NotificationsClient pb.SiteNotificationsClient
}

func RegisterRoutes(db *gorm.DB, pathToMediaDir string, notificationsClient pb.SiteNotificationsClient, secretKey string, r *gin.Engine) {
	h := handler{
		DB:                  db,
		PathToMediaDir:      pathToMediaDir,
		NotificationsClient: notificationsClient,
	}

	users := r.Group("/api/users")
	{
		users.GET("/", h.GetUsers)
		users.GET("/:id", h.GetUser)
		users.GET("/:id/profile-picture", h.GetUserProfilePicture)

		me := users.Group("/me")
		me.Use(middlewares.Auth(secretKey))
		{
			me.GET("/", h.GetMyProfile)
			me.GET("/profile-picture", h.GetMyProfilePicture)
			me.POST("/edited", h.EditProfile)
			me.PATCH("/on-verification", h.EditProfileOnVerification)
			me.PATCH("/settings", h.ChangeProfileSettings)
			me.DELETE("/", h.DeleteProfile)
		}
	}
}

func NewHandler(db *gorm.DB, notificationsClient pb.SiteNotificationsClient) handler {
	return handler{
		DB:                  db,
		NotificationsClient: notificationsClient,
	}
}
