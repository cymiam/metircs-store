package service

import (
	"github.com/cymiam/metrics-store/internal/repository"
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
}

func NewMetricService(config MetricServiceParams) *MetricService {
	return &MetricService{
		store:  config.Store,
		saver:  config.Saver,
		logger: config.Logger,
	}
}

func (service *MetricService) UpdateCounter(name string, newValue int64) {
	value, ok := service.store.GetCounter(name)
	if !ok {
		service.store.SetCounter(name, newValue)
		return
	}
	last := value[len(value)-1]
	service.store.SetCounter(name, last+newValue)
	if service.saver != nil {
		service.saver.OnMetricChanged()
	}
}

func (service *MetricService) UpdateGauge(name string, value float64) {
	service.store.SetGauge(name, value)
	if service.saver != nil {
		service.saver.OnMetricChanged()
	}
}

func (service *MetricService) GetCounter(name string) ([]int64, bool) {
	return service.store.GetCounter(name)
}

func (service *MetricService) GetGauge(name string) (float64, bool) {
	return service.store.GetGauge(name)
}

func (service *MetricService) GetCounters() map[string][]int64 {
	return service.store.GetCounters()
}

func (service *MetricService) GetGauges() map[string]float64 {
	return service.store.GetGauges()
}
