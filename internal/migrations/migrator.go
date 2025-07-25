package migrations

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	slog "github.com/salamsites/package-log"
)

func RunMigrations(logger *slog.Logger, db *sql.DB) error {
	err := goose.Up(db, "db/migrations")
	if err != nil {
		logger.Errorf("Migration failed: %v", err)
		return err
	}

	logger.Info("✅ Migrations applied successfully")
	return nil
}
