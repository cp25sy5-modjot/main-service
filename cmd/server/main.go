package main

import (
	"github.com/cp25sy5-modjot/main-service/internal/config"
	"github.com/cp25sy5-modjot/main-service/internal/database"
	server "github.com/cp25sy5-modjot/main-service/internal/httpapi"
)

func main() {
	conf := config.LoadConfig()
	db := database.NewPostgresDatabase(conf)
	database.AutoMigrate(db.GetDb())
	server.NewFiberServer(conf, db).Start()
}
