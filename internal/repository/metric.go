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
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string][]int64 `json:"counters"`
}

func NewStore() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string][]int64),
	}
}

func (m *MemStorage) SetCounter(name string, value int64) {
	m.Counters[name] = append(m.Counters[name], value)
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.Gauges[name] = value
}

func (m *MemStorage) GetCounter(name string) ([]int64, bool) {
	value, ok := m.Counters[name]
	return value, ok
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	value, ok := m.Gauges[name]
	return value, ok
}

func (m *MemStorage) GetGauges() map[string]float64 {
	cpy := make(map[string]float64, len(m.Gauges))

	for key, value := range m.Gauges {
		cpy[key] = value
	}

	return cpy
}

func (m *MemStorage) GetCounters() map[string][]int64 {
	cpy := make(map[string][]int64, len(m.Counters))

	for key, value := range m.Counters {
		tmp := make([]int64, len(value))
		copy(tmp, value)
		cpy[key] = tmp
	}

	return cpy
}
