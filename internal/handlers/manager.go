package handlers

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/handlers/http"
	"autotm-admin/internal/repository"
	"autotm-admin/internal/services"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	sminio "github.com/salamsites/minio-pkg"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
	spsql "github.com/salamsites/package-psql"
)

const (
	baseURL     = "/api/v1/autotm-admin"
	filesURL    = baseURL + "/files"
	brandURL    = baseURL + "/brand"
	settingsURL = baseURL + "/settings"
	regionsURL  = baseURL + "/regions"
	slidersURL  = baseURL + "/sliders"
	stocksURL   = baseURL + "/stocks"
	usersURL    = baseURL + "/users"
)

func Manager(logger *slog.Logger, clientPsql spsql.Client, minioFileClient sminio.ImageClient, cfg *configs.Config) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	newMiddleware := shttp.NewMiddleware(logger, cfg.Auth.JwtRegistration, nil)

	r.Route(filesURL, func(subRouter chi.Router) {
		filesHandler := http.NewFilesHandler(logger, newMiddleware, minioFileClient)
		filesHandler.FilesRegisterRoutes(subRouter)
	})

	r.Route(brandURL, func(subRouter chi.Router) {
		brandRepo := repository.NewBrandPsqlRepository(logger, clientPsql)
		brandService := services.NewBrandService(logger, brandRepo, minioFileClient)
		brandHandler := http.NewBrandHandler(logger, newMiddleware, brandService)
		brandHandler.BrandRegisterRoutes(subRouter)
	})

	r.Route(settingsURL, func(subRouter chi.Router) {
		settingsRepo := repository.NewSettingsPsqlRepository(logger, clientPsql)
		settingsService := services.NewSettingsService(logger, settingsRepo)

		if err := settingsService.InitSuperAdmin(context.Background()); err != nil {
			logger.Errorf("Failed to initialize super admin settings: %v", err)
		}
		settingsHandler := http.NewSettingsHandler(logger, newMiddleware, settingsService)
		settingsHandler.SettingsRegisterRoutes(subRouter)
	})

	r.Route(regionsURL, func(subRouter chi.Router) {
		regionsRepo := repository.NewRegionsPsqlRepository(logger, clientPsql)
		regionsService := services.NewRegionsService(logger, regionsRepo)
		regionsHandler := http.NewRegionsHandler(logger, newMiddleware, regionsService)
		regionsHandler.RegionsRegisterRoutes(subRouter)
	})

	r.Route(slidersURL, func(subRouter chi.Router) {
		sliderRepo := repository.NewSliderPsqlRepository(logger, clientPsql)
		sliderService := services.NewSlidersService(logger, sliderRepo, minioFileClient)
		sliderHandler := http.NewSliderHandler(logger, newMiddleware, sliderService)
		sliderHandler.SliderRegisterRoutes(subRouter)
	})

	r.Route(stocksURL, func(subRouter chi.Router) {
		stockRepo := repository.NewStockPsqlRepository(logger, clientPsql)
		stockService := services.NewStockService(logger, stockRepo)
		stockHandler := http.NewStockHandler(logger, newMiddleware, stockService)
		stockHandler.StockRegisterRoutes(subRouter)
	})

	r.Route(usersURL, func(subRouter chi.Router) {
		usersRepo := repository.NewUserPsqlRepository(logger, clientPsql)
		usersService := services.NewUserService(logger, usersRepo)
		usersHandler := http.NewUsersHandler(logger, newMiddleware, usersService)
		usersHandler.UsersRegisterRoutes(subRouter)
	})

	return r
}
