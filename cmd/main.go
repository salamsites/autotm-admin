package main

import (
	_ "autotm-admin/docs"
	"autotm-admin/internal/configs"
	"autotm-admin/internal/handlers"
	"autotm-admin/internal/migrations"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	sminio "github.com/salamsites/minio-pkg"
	files "github.com/salamsites/minio-pkg/client/files"
	file "github.com/salamsites/minio-pkg/client/image"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title AutoTM-Admin
// @version 3.0.0
// @description AutoTM-Admin swagger
// @externalDocs.description  Other services
// @BasePath /api/v1/autotm-admin
// @schemes http https
// @securityDefinitions.apiKey  ApiKeyAuth
// @security
// @in header
// @name authorization
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := configs.GetConfig()
	logger := slog.GetLogger(cfg.Log.Path, cfg.Log.Filename)

	//postgresql
	psqlClient, err := spsql.NewClient(ctx,
		spsql.Options{
			Host:          cfg.Storage.Psql.Host,
			Port:          cfg.Storage.Psql.Port,
			Database:      cfg.Storage.Psql.Database,
			Username:      cfg.Storage.Psql.Username,
			Password:      cfg.Storage.Psql.Password,
			PgPoolMaxConn: cfg.Storage.Psql.PgPoolMaxConn,
		})

	if err != nil {
		logger.Error("psql client does not connect: %v", err)
		panic("psql client does not connect")
	}
	defer psqlClient.Close()
	logger.Info("psql client connected")

	// Run db migrations
	if cfg.Storage.Psql.Migration {
		if err = migrations.RunMigrations(logger, psqlClient.StdDB()); err != nil {
			logger.Fatalf("Failed to apply migrations: %v", err)
		}
	} else {
		logger.Info("Migration flag is disabled; skipping migrations.")
	}

	minioImageClient, errImage := file.NewImageClient(sminio.Options{
		Endpoint:        "10.192.1.127:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	})

	if errImage != nil {
		logger.Fatal(errImage)
	}

	minioFileClient, errFile := files.NewFileClient(sminio.Options{
		Endpoint:        "10.192.1.127:9000",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	})

	if errFile != nil {
		logger.Fatal(errImage)
	}

	router := handlers.Manager(logger, psqlClient, minioImageClient, minioFileClient, cfg)

	router.Get("/autotm-admin/swagger/*", httpSwagger.WrapHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
		AllowedHeaders: []string{
			"*",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(router)

	srv := &http.Server{
		Addr:              cfg.Listen.Port,
		Handler:           handler,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	// Start the server in a separate Goroutine
	go func() {
		logger.Info("Starting the server on port", cfg.Listen.Port)
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()

	logger.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("final")
}
