package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"sync"

	"github.com/cp25sy5-modjot/main-service/internal/config"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn, err := utils.PostgresUrlBuilder(conf)
		if err != nil {
			panic("failed to build database URL")
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
