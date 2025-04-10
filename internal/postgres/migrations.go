package postgres

import (
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

//go:embed sql/migrations/*.sql
var migrationFiles embed.FS

func Migrate(log *slog.Logger, db *sqlx.DB) error {
	log.Debug("running migration")
	files, err := iofs.New(migrationFiles, "migrations") // get migrations from
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", files, "pgx", driver)
	if err != nil {
		return err
	}

	err = m.Up()

	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Error("migration failed", "error", err)
			return err
		}
		log.Debug("migration did not change anything")
	}

	log.Debug("migration finished")
	return nil
}
