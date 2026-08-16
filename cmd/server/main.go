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

	saver := service.NewMetricSaver(service.MetricSaverParams{
		Path:          config.FileStoragePath,
		StoreInterval: config.StoreInterval,
		Store:         repository,
	})

	if config.Restore {
		err := saver.PopulateStore()
		if err != nil {
			logger.Log.Fatal("Cannot populate store", zap.Error(err))
		}
	}

	service := service.NewMetricService(service.MetricServiceParams{Store: repository, Saver: saver})
	metricHandler := handler.NewMetricHandler(service)
	r := handler.NewMetricRouter(metricHandler)

	if config.StoreInterval > 0 {
		go saver.Update()
	}

	logger.Log.Info("Running server", zap.String("address", config.Addr))
	log.Fatal(http.ListenAndServe(config.Addr, r))

}
