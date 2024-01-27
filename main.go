package main

import (
	appuser "eau-de-go/internal/app_user"
	"eau-de-go/internal/db"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/transport/http"
	log "github.com/sirupsen/logrus"
)

func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting Up Our APP")

	database, err := db.NewDatabase()
	if err != nil {
		log.Error("failed to setup connection to the database")
		return err
	}

	// database.Ping(context.Background())

	queries := repository.New(database.Client)
	appUserService := appuser.NewAppUserService(queries)
	handler := http.NewHandler(appUserService)

	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
