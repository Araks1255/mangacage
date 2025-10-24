package testenv

import (
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Env struct {
	DB                  *gorm.DB
	MongoDB             *mongo.Database
	NotificationsClient pb.SiteNotificationsClient
	SecretKey           string
}
