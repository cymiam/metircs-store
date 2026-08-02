package repository

type MetricRepository interface {
	SetCounter(name string, value int64)
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, bool)
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewStore() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) SetCounter(name string, value int64) {
	m.counters[name] = value
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	value, ok := m.counters[name]
	return value, ok
}
