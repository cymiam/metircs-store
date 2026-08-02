package service

import "github.com/cymiam/metircs-store/internal/repository"

type MetricService struct {
	store repository.MetricRepository
}

func NewMetricService() *MetricService {
	return &MetricService{
		store: repository.NewStore(),
	}
}

func (service *MetricService) UpdateCounter(name string) {
	value, ok := service.store.GetCounter(name)
	if !ok {
		service.store.SetCounter(name, 1)
		return
	}
	service.store.SetCounter(name, value+1)
}

func (service *MetricService) UpdateGauge(name string, value float64) {
	service.store.SetGauge(name, value)
}

func (service *MetricService) GetCounter(name string) (int64, bool) {
	return service.store.GetCounter(name)
}

func (service *MetricService) GetGauge(name string) (float64, bool) {
	return service.store.GetGauge(name)
}

func (service *MetricService) GetCounters() map[string]int64 {
	return service.store.GetCounters()
}

func (service *MetricService) GetGauges() map[string]float64 {
	return service.store.GetGauges()
}
