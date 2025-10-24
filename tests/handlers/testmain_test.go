package handlers

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/internal/migrations"
	"github.com/Araks1255/mangacage/internal/seeder"
	dbPackage "github.com/Araks1255/mangacage/pkg/common/db"
	"github.com/Araks1255/mangacage/tests/testenv"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

var env testenv.Env

func TestMain(m *testing.M) {
	os.Chdir("./../..")
	viper.SetConfigFile("./pkg/common/envs/.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	dbUrl := viper.Get("DB_TEST_URL").(string)
	mongoUrl := viper.Get("MONGO_URL").(string)
	secretKey := viper.Get("SECRET_KEY").(string)

	db, err := dbPackage.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	mongoClient, err := dbPackage.MongoInit(ctx, mongoUrl)
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	mongoDB := mongoClient.Database("mangacage_test")

	env.DB = db
	env.MongoDB = mongoDB
	env.SecretKey = secretKey
	env.NotificationsClient = pb.NewSiteNotificationsClient(conn)

	migrateFlag := flag.Bool("migrate", false, "Run migrations with api")

	flag.Parse()

	if *migrateFlag {
		if err = migrations.Migrate(ctx, db, mongoDB); err != nil {
			panic(err)
		}
	}

	if err = seeder.Seed(db, mongoDB, "prod"); err != nil {
		panic(err)
	}

	code := m.Run()

	cleanTestDB(env.DB, env.MongoDB)

	os.Exit(code)
}

func cleanTestDB(db *gorm.DB, mongoDB *mongo.Database) {
	db.Exec("TRUNCATE TABLE authors, chapters, titles, users, teams, genres, tags RESTART IDENTITY CASCADE")

	ctx := context.Background()
	coll := mongoDB.Collection

	coll(mongodb.ChaptersPagesCollection).DeleteMany(ctx, bson.M{})
	coll(mongodb.TeamsCoversCollection).DeleteMany(ctx, bson.M{})
	coll(mongodb.TitlesCoversCollection).DeleteMany(ctx, bson.M{})
	coll(mongodb.UsersProfilePicturesCollection).DeleteMany(ctx, bson.M{})
}
