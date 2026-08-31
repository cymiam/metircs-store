package repository

import (
	"context"
	"fmt"

	models "github.com/cymiam/metrics-store/internal/model"
)

type MemStorage struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

func NewStore() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *MemStorage) GetAll(ctx context.Context) ([]models.Metric, error) {
	metrics := make([]models.Metric, 0, len(m.Gauges)+len(m.Counters))

	for name, value := range m.Gauges {
		gaugeValue := value

		metrics = append(metrics, models.Metric{
			ID:    name,
			MType: "gauge",
			Value: &gaugeValue,
		})
	}

	for name, value := range m.Counters {
		counterDelta := value

		metrics = append(metrics, models.Metric{
			ID:    name,
			MType: "counter",
			Delta: &counterDelta,
		})

	}

	return metrics, nil
}

func (m *MemStorage) GetMetric(ctx context.Context, name, metricType string) (models.Metric, error) {

	switch metricType {
	case "gauge":
		value, ok := m.Gauges[name]
		if !ok {
			return models.Metric{}, fmt.Errorf("Metric(%s):%s not found\n", metricType, name)
		}

		return models.Metric{
			ID:    name,
			MType: metricType,
			Value: &value,
		}, nil
	case "counter":
		value, ok := m.Counters[name]
		if !ok {
			return models.Metric{}, fmt.Errorf("Metric(%s):%s not found\n", metricType, name)
		}
		return models.Metric{
			ID:    name,
			MType: metricType,
			Delta: &value,
		}, nil
	}
	return models.Metric{}, fmt.Errorf("Unknown metric type, %s\n", metricType)
}

func (m *MemStorage) SetMetric(ctx context.Context, metric models.Metric) error {
	switch metric.MType {
	case "gauge":
		if metric.Value == nil {
			return fmt.Errorf("%s value is nil", metric)
		}
		m.Gauges[metric.ID] = *metric.Value
		return nil
	case "counter":
		if metric.Delta == nil {
			return fmt.Errorf("%s value is nil", metric)
		}
		m.Counters[metric.ID] += *metric.Delta
		return nil
	}
	return fmt.Errorf("Unknown metric type, %s\n", metric.MType)
}
