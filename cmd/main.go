package main

import (
	"github.com/Araks1255/mangabrad/pkg/common/db"
	"github.com/Araks1255/mangabrad/pkg/handlers/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := viper.Get("DB_URL").(string)

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	auth.RegisterRoutes(db, router)

	router.Run(":8080")
}
