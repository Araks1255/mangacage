package handlers

import (
	"flag"
	"os"
	"testing"

	"github.com/Araks1255/mangacage/internal/migrations"
	"github.com/Araks1255/mangacage/internal/seeder"
	dbPackage "github.com/Araks1255/mangacage/pkg/common/db"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB // Я не нашёл другого нормального способа организовать доступ к бд в тестах. Выбор был между глобальными переменными и инициализацией отдельного подключения в каждом тесте
	mongoDB *mongo.Database
)

func TestMain(m *testing.M) {
	os.Chdir("./../..")
	viper.SetConfigFile("./pkg/common/envs/.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	dbUrl := viper.Get("DB_TEST_URL").(string)

	var err error
	db, err = dbPackage.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	mongoUrl := viper.Get("MONGO_URL").(string)

	mongoClient, err := dbPackage.MongoInit(mongoUrl)
	if err != nil {
		panic(err)
	}

	mongoDB = mongoClient.Database("mangacage_test")

	migrateFlag := flag.Bool("migrate", false, "Run migrations with api")

	seedMode := flag.String("seed", "", "Mode of seed")

	flag.Parse()

	if *migrateFlag {
		if err = migrations.GormMigrate(db); err != nil {
			panic(err)
		}
	}

	if *seedMode != "" {
		if err = seeder.Seed(db, mongoDB, *seedMode); err != nil {
			panic(err)
		}
	}

	code := m.Run()

	db.Exec("TRUNCATE TABLE authors, chapters, titles, users, teams, volumes RESTART IDENTITY CASCADE") // В конце всех тестов бд очищается. Иначе там с транзакциями бы запара была. Так что запуск тестов должен сопровождаться флагом --seed=test

	os.Exit(code)
}
