package testenv

import (
	pb "github.com/Araks1255/mangacage_protos"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Env struct {
	DB                  *gorm.DB
	MongoDB             *mongo.Database
	NotificationsClient pb.NotificationsClient
	SecretKey           string
}
