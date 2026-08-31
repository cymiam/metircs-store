package repository

import (
	"context"

	models "github.com/cymiam/metrics-store/internal/model"
)

type MetricRepository interface {
	SetMetric(ctx context.Context, metirc models.Metric) error
	GetMetric(ctx context.Context, name, metricType string) (models.Metric, error)
	GetAll(ctx context.Context) ([]models.Metric, error)
}
