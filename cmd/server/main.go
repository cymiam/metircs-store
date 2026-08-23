package main

import (
	"log"
	"net/http"

	config "github.com/cymiam/metrics-store/internal/config/server"
	"github.com/cymiam/metrics-store/internal/handler"
	"github.com/cymiam/metrics-store/internal/logger"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/cymiam/metrics-store/internal/service"
	"go.uber.org/zap"
)

func main() {

	baseLog, err := logger.NewLogger("info", "main")
	if err != nil {
		log.Fatal("Cannot start logger: ", err)
	}
	defer func() {
		_ = baseLog.Sync()
	}()

	handlerLog := baseLog.With(zap.String("layer", "handler"))
	httpLog := baseLog.With(zap.String("layer", "http"))
	saverLog := baseLog.With(zap.String("layer", "service"))

	config, err := config.ParseServerConfig()

	repository := repository.NewStore()

	saver, err := service.NewMetricSaver(service.MetricSaverParams{
		Path:          config.FileStoragePath,
		StoreInterval: config.StoreInterval,
		Store:         repository,
		Logger:        saverLog,
	})

	if err != nil {
		log.Fatal("Metric Saver error", err)
	}

	if config.Restore {
		err := saver.PopulateStore()
		if err != nil {
			log.Fatal("Cannot populate store", err)
		}
	}

	service := service.NewMetricService(service.MetricServiceParams{Store: repository, Saver: saver, Logger: saverLog})
	metricHandler := handler.NewMetricHandler(service, handlerLog)
	r := handler.NewMetricRouter(metricHandler, httpLog)

	if config.StoreInterval > 0 {
		go saver.StartTicker()
	}

	baseLog.Info("Running server", zap.String("address", config.Addr))
	log.Fatal(http.ListenAndServe(config.Addr, r))

}
