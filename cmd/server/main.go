package main

import (
	"modjot/internal/config"
	"modjot/internal/database"
	server "modjot/internal/httpapi"
)

func main() {
	conf := config.LoadConfig()
	db := database.NewPostgresDatabase(conf)
	database.AutoMigrate(db.GetDb())
	server.NewFiberServer(conf, db).Start()
}
