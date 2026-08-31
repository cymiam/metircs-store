package service

import (
	"context"
	"time"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type MetricService struct {
	store  repository.MetricRepository
	saver  *MetricSaver
	logger *zap.Logger
	db     *pgx.Conn
}

type MetricServiceParams struct {
	Store  repository.MetricRepository
	Saver  *MetricSaver
	Logger *zap.Logger
	DB     *pgx.Conn
}

func NewMetricService(config MetricServiceParams) *MetricService {
	return &MetricService{
		store:  config.Store,
		saver:  config.Saver,
		logger: config.Logger,
		db:     config.DB,
	}
}

func (service *MetricService) UpdateCounter(name string, delta int64) {

	service.store.SetMetric(context.TODO(), models.Metric{ID: name, MType: "counter", Delta: &delta})

	if service.saver != nil {
		metric := models.Metric{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		}
		service.saver.OnMetricChanged(metric)
	}
}

func (service *MetricService) UpdateGauge(name string, value float64) {
	service.store.SetMetric(context.TODO(), models.Metric{ID: name, MType: "gauge", Value: &value})
	if service.saver != nil {
		metric := models.Metric{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		service.saver.OnMetricChanged(metric)
	}
}

func (service *MetricService) GetMetric(name, metricType string) (models.Metric, error) {
	return service.store.GetMetric(context.TODO(), name, metricType)
}

func (service *MetricService) GetAll() ([]models.Metric, error) {
	return service.store.GetAll(context.TODO())
}

func (service *MetricService) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return service.db.Ping(ctx)
}
