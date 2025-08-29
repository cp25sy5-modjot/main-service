package main

import (
	"modjot/config"
	"modjot/database"
	"modjot/server"
)

func main() {
	conf := config.LoadConfig()
	db := database.NewPostgresDatabase(conf)
	database.AutoMigrate(db.GetDb())
	server.NewFiberServer(conf, db).Start()
}
