package handlers

import (
	"autotm-admin/internal/configs"
	"autotm-admin/internal/handlers/http"
	"autotm-admin/internal/helpers"
	"autotm-admin/internal/repository"
	"autotm-admin/internal/services"
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	es "github.com/elastic/go-elasticsearch/v8"
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
	carsURL     = baseURL + "/cars"
)

func Manager(logger *slog.Logger, clientPsql spsql.Client, minioImageClient sminio.ImageClient, minioFileClient sminio.FileClient, cfg *configs.Config, esClient *es.Client) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	newMiddleware := shttp.NewMiddleware(logger, cfg.Auth.JwtRegistration, nil)

	pushService := services.NewPushService(logger, cfg)
	repo := repository.NewUserPsqlRepository(logger, clientPsql)
	userService := services.NewUserService(logger, repo)

	r.Route(filesURL, func(subRouter chi.Router) {
		filesHandler := http.NewFilesHandler(logger, newMiddleware, minioImageClient)
		filesHandler.FilesRegisterRoutes(subRouter)
	})

	r.Route(brandURL, func(subRouter chi.Router) {
		brandRepo := repository.NewBrandPsqlRepository(logger, clientPsql)
		brandService := services.NewBrandService(logger, brandRepo, minioImageClient)
		brandHandler := http.NewBrandHandler(logger, newMiddleware, brandService)
		brandHandler.BrandRegisterRoutes(subRouter)
	})

	r.Route(settingsURL, func(subRouter chi.Router) {
		settingsRepo := repository.NewSettingsPsqlRepository(logger, clientPsql)
		settingsService := services.NewSettingsService(logger, settingsRepo, cfg)

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
		sliderService := services.NewSlidersService(logger, sliderRepo, minioImageClient)
		sliderHandler := http.NewSliderHandler(logger, newMiddleware, sliderService)
		sliderHandler.SliderRegisterRoutes(subRouter)
	})

	r.Route(stocksURL, func(subRouter chi.Router) {
		stockRepo := repository.NewStockPsqlRepository(logger, clientPsql, trmpgx.DefaultCtxGetter)
		stockService := services.NewStockService(logger, stockRepo, userService, pushService)
		stockHandler := http.NewStockHandler(logger, newMiddleware, clientPsql, stockService, minioFileClient, minioImageClient)
		stockHandler.StockRegisterRoutes(subRouter)
	})

	r.Route(usersURL, func(subRouter chi.Router) {
		usersRepo := repository.NewUserPsqlRepository(logger, clientPsql)
		usersService := services.NewUserService(logger, usersRepo)
		usersHandler := http.NewUsersHandler(logger, newMiddleware, usersService)
		usersHandler.UsersRegisterRoutes(subRouter)
	})

	r.Route(carsURL, func(subRouter chi.Router) {
		carsRepo := repository.NewCarsPsqlRepository(logger, clientPsql)
		stockRepo := repository.NewStockPsqlRepository(logger, clientPsql, trmpgx.DefaultCtxGetter)
		carsService := services.NewCarsService(logger, carsRepo, userService, pushService, stockRepo, esClient, helpers.CarIndexName)
		carsHandler := http.NewCarsHandler(logger, newMiddleware, carsService)
		carsHandler.CarsRegisterRoutes(subRouter)
	})

	return r
}
