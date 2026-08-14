package main

import (
	"log"
	"net/http"

	config "github.com/cymiam/metircs-store/internal/config/server"
	"github.com/cymiam/metircs-store/internal/handler"
	"github.com/cymiam/metircs-store/internal/logger"
	"github.com/cymiam/metircs-store/internal/repository"
	"github.com/cymiam/metircs-store/internal/service"
	"go.uber.org/zap"
)

func main() {
	if err := logger.Initialize("Info"); err != nil {
		panic(err)
	}
	config := config.ParseServerConfig()

	repository := repository.NewStore()
	service := service.NewMetricService(repository)
	metricHandler := handler.NewMetricHandler(service)

	r := handler.NewMetricRouter(metricHandler)
	logger.Log.Info("Running server", zap.String("address", config.Addr))
	log.Fatal(http.ListenAndServe(config.Addr, r))
}
