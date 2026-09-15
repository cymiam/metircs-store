package main

import (
	"context"
	"log"
	"net/http"
	"time"

	config "github.com/cymiam/metrics-store/internal/config/server"
	"github.com/cymiam/metrics-store/internal/handler"
	"github.com/cymiam/metrics-store/internal/logger"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/cymiam/metrics-store/internal/service"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
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

	if err != nil {
		log.Fatal("Cannot parse server config: ", err)
	}

	var metricRepository repository.MetricRepository
	var saver *service.MetricSaver
	var pool *pgxpool.Pool

	if config.ConnectionString == "" {

		baseLog.Info("Run server with file and memory storage")

		store := repository.NewStore()

		metricRepository = store

		saver, err = service.NewMetricSaver(service.MetricSaverParams{
			Path:          config.FileStoragePath,
			StoreInterval: config.StoreInterval,
			Store:         store,
			Logger:        saverLog,
			Restore:       config.Restore,
		})

		if err != nil {
			log.Fatal("Metric Saver error", err)
		}

		if config.StoreInterval > 0 {
			go saver.StartTicker()
		}

	} else {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		pool, err = pgxpool.New(ctx, config.ConnectionString)

		if err != nil {
			log.Fatal("Unable to connect to database: %w", err)
		}
		defer pool.Close()

		baseLog.Info("Run server with postgres storage")
		baseLog.Info("Running migrations")

		err = runMigrations(pool)

		if err != nil {
			log.Fatal("Error running migrations: %w", err)
		}

		baseLog.Info("Migratios run succsess")

		metricRepository = repository.NewPostgresStorage(repository.PostgreStorageParams{Pool: pool})
	}

	metricService := service.NewMetricService(service.MetricServiceParams{Store: metricRepository, Saver: saver, Logger: saverLog})
	metricHandler := handler.NewMetricHandler(metricService, handlerLog)

	healthService := service.NewHealthService(pool)
	healthHadler := handler.NewHealthHandler(handler.HealthHandlerParams{HealthService: healthService, Logger: httpLog})

	mainRouter := handler.NewRouter(handler.MainRouterParams{Logger: httpLog, MetricHandler: metricHandler, HealthHandler: healthHadler, Key: config.Key})

	baseLog.Info("Running server", zap.String("address", config.Addr))
	log.Fatal(http.ListenAndServe(config.Addr, mainRouter))

}

func runMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
