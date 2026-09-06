package service

import (
	"context"
	"fmt"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/repository"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type MetricService struct {
	store  repository.MetricRepository
	saver  *MetricSaver
	logger *zap.Logger
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
	}
}

func (service *MetricService) UpdateCounter(name string, delta int64) error {

	err := service.store.SetMetric(context.TODO(), models.Metric{ID: name, MType: "counter", Delta: &delta})

	if err != nil {
		return fmt.Errorf("update counter: %w", err)
	}

	if service.saver != nil {
		metric := models.Metric{
			ID:    name,
			MType: "counter",
			Delta: &delta,
		}
		service.saver.OnMetricChanged(metric)
	}

	return nil
}

func (service *MetricService) UpdateGauge(name string, value float64) error {
	err := service.store.SetMetric(context.TODO(), models.Metric{ID: name, MType: "gauge", Value: &value})

	if err != nil {
		return fmt.Errorf("update gauge: %w", err)
	}
	if service.saver != nil {
		metric := models.Metric{
			ID:    name,
			MType: "gauge",
			Value: &value,
		}
		service.saver.OnMetricChanged(metric)
	}

	return nil
}

func (service *MetricService) GetMetric(name, metricType string) (models.Metric, error) {
	return service.store.GetMetric(context.TODO(), name, metricType)
}

func (service *MetricService) GetAll() ([]models.Metric, error) {
	return service.store.GetAll(context.TODO())
}
