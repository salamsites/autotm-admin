package main

import (
	_ "autotm-admin/docs"
	"autotm-admin/internal/configs"
	"autotm-admin/internal/handlers"
	"context"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// @title AutoTM-Admin
// @version 3.0.0
// @description AutoTM-Admin swagger
// @externalDocs.description  Other services
// @BasePath /autotm-admin
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

	router := handlers.Manager(logger, psqlClient, cfg)

	router.Get("/autotm-admin/swagger/*", httpSwagger.WrapHandler)

	fileServer := http.FileServer(http.Dir(cfg.FilePath))
	router.Handle("/uploads/*", http.StripPrefix("/uploads", fileServer))

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
