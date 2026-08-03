package repository

type MetricRepository interface {
	SetCounter(name string, value int64)
	SetGauge(name string, value float64)
	GetCounter(name string) ([]int64, bool)
	GetGauge(name string) (float64, bool)
	GetGauges() map[string]float64
	GetCounters() map[string][]int64
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string][]int64
}

func NewStore() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string][]int64),
	}
}

func (m *MemStorage) SetCounter(name string, value int64) {
	m.counters[name] = append(m.counters[name], value)
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.gauges[name] = value
}

func (m *MemStorage) GetCounter(name string) ([]int64, bool) {
	value, ok := m.counters[name]
	return value, ok
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := m.gauges[name]
	return value, ok
}

func (m *MemStorage) GetGauges() map[string]float64 {
	cpy := make(map[string]float64, len(m.gauges))

	for key, value := range m.gauges {
		cpy[key] = value
	}

	return cpy
}

func (m *MemStorage) GetCounters() map[string][]int64 {
	cpy := make(map[string][]int64, len(m.counters))

	for key, value := range m.counters {
		tmp := make([]int64, len(value))
		copy(tmp, value)
		cpy[key] = tmp
	}

	return cpy
}
