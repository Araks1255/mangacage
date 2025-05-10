package main

import (
	"context"
	"flag"

	"github.com/Araks1255/mangacage/internal/migrations"
	"github.com/Araks1255/mangacage/internal/seeder"
	"github.com/Araks1255/mangacage/pkg/common/db"
	"github.com/Araks1255/mangacage/pkg/handlers/auth"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/handlers/search"
	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/pkg/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/handlers/users/favorites"
	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/handlers/views"
	"github.com/Araks1255/mangacage/pkg/handlers/volumes"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := viper.Get("SECRET_KEY").(string)
	dbUrl := viper.Get("DB_URL").(string)
	mongoUrl := viper.Get("MONGO_URL").(string)

	mongoClient, err := db.MongoInit(mongoUrl)
	if err != nil {
		mongoClient.Disconnect(context.TODO())
		panic(err)
	}
	defer mongoClient.Disconnect(context.TODO())

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	notificationsClient := pb.NewNotificationsClient(conn)

	migrateFlag := flag.Bool("migrate", false, "Run migrations with api") // Получение cli флага. Если будет запуск: go run cmd/main.go --migrate, то запустится миграция бд
	seedMode := flag.String("seed", "", "Mode of seed")

	flag.Parse()

	if *migrateFlag {
		if err := migrations.GormMigrate(db); err != nil {
			panic(err)
		}
	}

	if *seedMode != "" {
		if err = seeder.Seed(db, mongoClient.Database("mangacage"), *seedMode); err != nil {
			panic(err)
		}
	}

	router := gin.Default()

	auth.RegisterRoutes(db, mongoClient, notificationsClient, secretKey, router)
	titles.RegisterRoutes(db, mongoClient, notificationsClient, router)
	teams.RegisterRoutes(db, mongoClient, secretKey, router)
	joinrequests.RegisterRoutes(db, secretKey, router)
	participants.RegisterRoutes(db, secretKey, router)
	chapters.RegisterRoutes(db, mongoClient, notificationsClient, secretKey, router)
	volumes.RegisterRoutes(db, notificationsClient, secretKey, router)
	search.RegisterRoutes(db, router)
	users.RegisterRoutes(db, mongoClient, notificationsClient, secretKey, router)
	views.RegisterRoutes(db, router)
	favorites.RegisterRoutes(db, secretKey, router)
	moderation.RegisterRoutes(db, mongoClient, secretKey, router)

	router.Run(":8080")
}
