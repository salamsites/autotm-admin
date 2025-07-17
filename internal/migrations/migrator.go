package migrations

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	slog "github.com/salamsites/package-log"
)

func RunMigrations(logger *slog.Logger, dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Errorf("Failed to open DB: %v", err)
		return err
	}
	defer db.Close()

	logger.Println("🚀 Running goose migrations...")

	if err = goose.SetDialect("postgres"); err != nil {
		logger.Errorf("Failed to set goose dialect: %v", err)
		return err
	}

	err = goose.Up(db, "db/migrations")
	if err != nil {
		logger.Errorf("Migration failed: %v", err)
		return err
	}

	logger.Info("✅ Migrations applied successfully")
	return nil
}
