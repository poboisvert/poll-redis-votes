package db

import (
	"database/sql"
	"fmt"

	"github.com/poboisvert/poll-redis-pubsub/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var dbInstance *sql.DB

func Connect() error {
	var err error
	dbConfig := config.GetConfig()

	dbInstance, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database))
	if err != nil {
		return err
	}

	if err = dbInstance.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to PostgreSQL")
	return nil
}

func Close() {
	if dbInstance != nil {
		_ = dbInstance.Close()
	}
}

func GetDB() *sql.DB {
	return dbInstance
}
