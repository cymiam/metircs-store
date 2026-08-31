package repository

import (
	"context"

	models "github.com/cymiam/metrics-store/internal/model"
)

type PostrgreStorage struct {
}

// GetAll implements [MetricRepository].
func (p *PostrgreStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	panic("unimplemented")
}

// GetMetric implements [MetricRepository].
func (p *PostrgreStorage) GetMetric(ctx context.Context, name string, metricType string) (*models.Metric, error) {
	panic("unimplemented")
}

// SetMetric implements [MetricRepository].
func (p *PostrgreStorage) SetMetric(ctx context.Context, metirc models.Metric) error {
	panic("unimplemented")
}
