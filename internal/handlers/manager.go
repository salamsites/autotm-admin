package handlers

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/handlers/http"
	"autotm-admin/internal/repository"
	"autotm-admin/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

const (
	baseURL  = "/autotm-admin"
	brandURL = baseURL + "/brand"
)

func Manager(logger *slog.Logger, clientPsql spsql.Client, cfg *configs.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	newMiddleware := shttp.NewMiddleware(logger, cfg.Auth.JwtRegistration, nil)

	r.Route(brandURL, func(subRouter chi.Router) {
		brandRepo := repository.NewBrandPsqlRepository(logger, clientPsql)
		brandService := services.NewBrandService(logger, brandRepo)
		brandHandler := http.NewBrandHandler(logger, newMiddleware, brandService)
		brandHandler.BrandRegisterRoutes(subRouter)
	})

	return r
}
