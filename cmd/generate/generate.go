package main

import (
	"autotm-admin/internal/configs"
	shttp "github.com/salamsites/package-http"
	slog "github.com/salamsites/package-log"
)

func main() {
	cfg := configs.GetConfig()
	logger := slog.GetLogger(cfg.Log.Path, cfg.Log.Filename)

	shttp.InitSwagger(logger, cfg.Swagger)
}
