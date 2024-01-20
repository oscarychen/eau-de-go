package main

import (
	"eau-de-go/internal/db"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func Migrate() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Running migration")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database, err := db.NewDatabase()
	if err != nil {
		log.Error("failed to setup connection to the database")
		return err
	}

	err = database.MigrateDB()
	if err != nil {
		log.Error("failed to migrate database")
		return err
	}

	return nil
}

func main() {
	if err := Migrate(); err != nil {
		log.Error(err)
		log.Fatal("Error running migration")
	}
}
