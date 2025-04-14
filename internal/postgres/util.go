package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"log/slog"

	"github.com/jmoiron/sqlx"
)

func SetupDBConnection(log *slog.Logger, address string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("Failed to connect to database", "address", address, "error", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping database", "error", err)
		return nil, err
	}

	return db, nil
}
