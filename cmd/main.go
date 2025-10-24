package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Araks1255/mangacage/internal/migrations"
	"github.com/Araks1255/mangacage/internal/seeder"
	cpc "github.com/Araks1255/mangacage/internal/workers/chapters_pages_compressor"
	rl "github.com/Araks1255/mangacage/pkg/logging/rotate_logger"
	"github.com/Araks1255/mangacage/pkg/common/db"
	"github.com/Araks1255/mangacage/pkg/handlers/auth"
	"github.com/Araks1255/mangacage/pkg/handlers/authors"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/handlers/genres"
	"github.com/Araks1255/mangacage/pkg/handlers/roles"
	"github.com/Araks1255/mangacage/pkg/handlers/tags"
	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/pkg/handlers/teams/joinrequests"
	"github.com/Araks1255/mangacage/pkg/handlers/teams/participants"
	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/handlers/titles/translaterequests"
	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/handlers/users/favorites"
	"github.com/Araks1255/mangacage/pkg/handlers/users/moderation"
	"github.com/Araks1255/mangacage/pkg/handlers/users/viewedchapters"
	"github.com/Araks1255/mangacage/pkg/handlers/views"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	ctx := context.Background()

	secretKey := viper.Get("SECRET_KEY").(string)
	dbUrl := viper.Get("DB_URL").(string)
	pathToMediaDir := viper.Get("PATH_TO_MEDIA_DIR").(string)
	pathToLogsDir := viper.Get("PATH_TO_LOGS_DIR").(string)

	rotateLogger, err := rl.NewRotateLogger(pathToLogsDir)
	if err != nil {
		panic(err)
	}

	log.SetOutput(rotateLogger)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	notificationsClient := pb.NewSiteNotificationsClient(conn)

	chaptersPagesCompressor := cpc.NewChaptersPagesCompressor(ctx, db, 16384)

	migrateFlag := flag.Bool("migrate", false, "Run migrations with api") // Получение cli флага. Если будет запуск: go run cmd/main.go --migrate, то запустится миграция бд
	seedMode := flag.String("seed", "", "Mode of seed")

	flag.Parse()

	if *migrateFlag {
		if err := migrations.Migrate(ctx, db, pathToMediaDir); err != nil {
			panic(err)
		}
	}

	if *seedMode != "" {
		if err = seeder.Seed(db, pathToMediaDir, *seedMode); err != nil {
			panic(err)
		}
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:80", "http://localhost", "http://localhost/"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.RegisterRoutes(db, notificationsClient, secretKey, router)
	titles.RegisterRoutes(db, pathToMediaDir, notificationsClient, secretKey, router)
	teams.RegisterRoutes(db, pathToMediaDir, notificationsClient, secretKey, router)
	joinrequests.RegisterRoutes(db, secretKey, notificationsClient, router)
	participants.RegisterRoutes(db, secretKey, notificationsClient, router)
	chapters.RegisterRoutes(db, pathToMediaDir, chaptersPagesCompressor, notificationsClient, secretKey, router)
	users.RegisterRoutes(db, pathToMediaDir, notificationsClient, secretKey, router)
	views.RegisterRoutes(db, router, secretKey)
	favorites.RegisterRoutes(db, secretKey, router)
	moderation.RegisterRoutes(db, secretKey, router)
	genres.RegisterRoutes(db, notificationsClient, secretKey, router)
	tags.RegisterRoutes(db, notificationsClient, secretKey, router)
	authors.RegisterRoutes(db, notificationsClient, secretKey, router)
	viewedchapters.RegisterRoutes(db, secretKey, router)
	roles.RegisterRoutes(db, router)
	translaterequests.RegisterRoutes(db, secretKey, notificationsClient, router)

	go chaptersPagesCompressor.Start()
	go rotateLogger.Start()

	router.Run(":8080")
}
