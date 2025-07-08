package migrations

import (
	"autotm-admin/internal/configs"
	"fmt"
	_ "github.com/lib/pq"
	slog "github.com/salamsites/package-log"
	"log"
	"os/exec"
)

func RunMigrations(logger *slog.Logger, cfg *configs.Config) {
	if !cfg.Storage.Psql.Migration {
		log.Println("ℹ️ Migration flag is disabled in config. Skipping migrations.")
		return
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.Psql.Username,
		cfg.Storage.Psql.Password,
		cfg.Storage.Psql.Host,
		cfg.Storage.Psql.Port,
		cfg.Storage.Psql.Database,
	)
	logger.Println("🚀 Running migrations...")

	cmd := exec.Command("migrate", "-path", "db/migrations", "-database", dsn, "up")

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Fatalf("❌ Migration failed: %v\n🧾 Details:\n%s", err, string(output))
	}

	logger.Info("✅ Migrations applied successfully:\n%s", string(output))
}
