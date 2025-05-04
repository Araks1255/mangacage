package chapters

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB                        *gorm.DB
	ChaptersOnModerationPages *mongo.Collection
	ChaptersPages             *mongo.Collection
	NotificationsClient       pb.NotificationsClient
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, notificationsClient pb.NotificationsClient, r *gin.Engine) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)

	chaptersOnModerationPagesCollection := client.Database("mangacage").Collection("chapters_on_moderation_pages")
	chapterPagesCollection := client.Database("mangacage").Collection("chapters_pages")

	h := handler{
		DB:                        db,
		ChaptersOnModerationPages: chaptersOnModerationPagesCollection,
		ChaptersPages:             chapterPagesCollection,
		NotificationsClient:       notificationsClient,
	}

	privateChapter := r.Group("/api/chapters")
	privateChapter.Use(middlewares.AuthMiddleware(secretKey))
	{
		privateChapter.POST("/", h.CreateChapter)
		privateChapter.DELETE("/:id", h.DeleteChapter)
		privateChapter.POST("/:id/edited", h.EditChapter) // Тут идёт создание отредактированной главы (прямо отдельная сущность в отдельной таблице базы данных), поэтому post а не put
	}

	publicChapter := r.Group("/api/chapters")
	{
		publicChapter.GET("/:id", h.GetChapter)
		publicChapter.GET("/:id/page/:page")
	}

	r.GET("/api/volume/:id/chapters", h.GetVolumeChapters)
}

func NewHandler(db *gorm.DB, notificationsClient pb.NotificationsClient, chaptersOnModerationPages, chaptersPages *mongo.Collection) handler {
	return handler{
		DB:                        db,
		ChaptersOnModerationPages: chaptersOnModerationPages,
		ChaptersPages:             chaptersPages,
		NotificationsClient:       notificationsClient,
	}
}
