package repository

import (
	"context"
	"fmt"
	"sync"

	models "github.com/cymiam/metrics-store/internal/model"
)

type MemStorage struct {
	mutex    sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewStore() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *MemStorage) GetAll(ctx context.Context) ([]models.Metric, error) {

	m.mutex.RLock()
	defer m.mutex.RUnlock()
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

	m.mutex.RLock()
	defer m.mutex.RUnlock()
	switch metricType {
	case "gauge":
		value, ok := m.Gauges[name]
		if !ok {
			return models.Metric{}, fmt.Errorf("Metric(%s):%s not found", metricType, name)
		}

		return models.Metric{
			ID:    name,
			MType: metricType,
			Value: &value,
		}, nil
	case "counter":
		value, ok := m.Counters[name]
		if !ok {
			return models.Metric{}, fmt.Errorf("Metric(%s):%s not found", metricType, name)
		}
		return models.Metric{
			ID:    name,
			MType: metricType,
			Delta: &value,
		}, nil
	}
	return models.Metric{}, fmt.Errorf("unknown metric type, %s", metricType)
}

func (m *MemStorage) SetMetric(ctx context.Context, metric models.Metric) error {

	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.setMetric(metric)

}

func (m *MemStorage) SetMetrics(ctx context.Context, metrics []models.Metric) error {

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, metric := range metrics {
		if err := m.setMetric(metric); err != nil {
			return fmt.Errorf("set batch metric: %w", err)
		}
	}

	return nil
}

func (m *MemStorage) setMetric(metric models.Metric) error {
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
	return fmt.Errorf("unknown metric type, %s", metric.MType)
}
