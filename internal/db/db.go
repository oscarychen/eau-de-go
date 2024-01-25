package db

import (
	"context"
	"eau-de-go/internal/settings"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	Client *sqlx.DB
}

// NewDatabase - returns a pointer to a database object
func NewDatabase() (*Database, error) {
	log.Info("Setting up new database connection")

	connectionString := fmt.Sprintf(
		"host=%s port=%s sslmode=%s user=%s dbname=%s password=%s",
		settings.DbHost,
		settings.DbPort,
		settings.SslMode,
		settings.DbUsername,
		settings.DbName,
		settings.DbPassword,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return &Database{}, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{
		Client: db,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
