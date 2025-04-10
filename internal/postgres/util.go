package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

func SetupDBConnection(log *slog.Logger, address, database, username, password string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, address, database)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database", "error", err)
		return nil, err
	}

	return db, nil
}
