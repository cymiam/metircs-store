package main

import (
	"context"
	"log"
	"net/http"

	config "github.com/cymiam/metrics-store/internal/config/server"
	"github.com/cymiam/metrics-store/internal/handler"
	"github.com/cymiam/metrics-store/internal/logger"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/cymiam/metrics-store/internal/service"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	conn, err := pgx.Connect(context.Background(), config.ConnectionString)

	if err != nil {
		log.Fatal("Unable to connect to database: %w\n", err)
	}
	defer conn.Close(context.Background())

	repository := repository.NewStore()

	saver, err := service.NewMetricSaver(service.MetricSaverParams{
		Path:          config.FileStoragePath,
		StoreInterval: config.StoreInterval,
		Store:         repository,
		Logger:        saverLog,
		Restore:       config.Restore,
	})

	if err != nil {
		log.Fatal("Metric Saver error", err)
	}

	service := service.NewMetricService(service.MetricServiceParams{Store: repository, Saver: saver, Logger: saverLog, DB: conn})
	metricHandler := handler.NewMetricHandler(service, handlerLog)
	r := handler.NewMetricRouter(metricHandler, httpLog)

	if config.StoreInterval > 0 {
		go saver.StartTicker()
	}

	baseLog.Info("Running server", zap.String("address", config.Addr))
	log.Fatal(http.ListenAndServe(config.Addr, r))

}
