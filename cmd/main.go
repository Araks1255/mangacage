package main

import (
	"context"

	"github.com/Araks1255/mangacage/pkg/common/db"
	"github.com/Araks1255/mangacage/pkg/handlers/auth"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/handlers/search"
	"github.com/Araks1255/mangacage/pkg/handlers/teams"
	"github.com/Araks1255/mangacage/pkg/handlers/titles"
	"github.com/Araks1255/mangacage/pkg/handlers/users"
	"github.com/Araks1255/mangacage/pkg/handlers/users/favorites"
	"github.com/Araks1255/mangacage/pkg/handlers/views"
	"github.com/Araks1255/mangacage/pkg/handlers/volumes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

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

	router := gin.Default()

	auth.RegisterRoutes(db, mongoClient, router)
	titles.RegisterRoutes(db, mongoClient, router)
	teams.RegisterRoutes(db, router)
	chapters.RegisterRoutes(db, mongoClient, router)
	volumes.RegisterRoutes(db, mongoClient, router)
	search.RegisterRoutes(db, router)
	users.RegisterRoutes(db, mongoClient, router)
	views.RegisterRoutes(db, router)
	favorites.RegisterRoutes(db, router)

	router.Run(":8080")
}
