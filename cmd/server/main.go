package main

import (
	"fmt"
	"net/http"
	"strconv"
)

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

var store = NewStore()

func (m *MemStorage) updateCounter(name string) {
	value, ok := m.counters[name]
	if !ok {
		m.counters[name] = 1
		return
	}
	m.counters[name] = value + 1
}

func (m *MemStorage) updateGauge(name string, value float64) {
	m.gauges[name] = value
}

func handleMetricUpdate(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("metric_type")
	metricName := r.PathValue("metric_name")
	metricValue, err := strconv.ParseFloat(r.PathValue("metric_value"), 64)
	fmt.Printf("Got request, %s, %s, %f\n", metricType, metricName, metricValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	switch metricType {
	case "counter":
		store.updateCounter(metricName)
		w.WriteHeader(http.StatusOK)
	case "gauge":
		store.updateGauge(metricName, metricValue)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{metric_type}/{metric_name}/{metric_value}", handleMetricUpdate)
	err := http.ListenAndServe("127.0.0.1:8080", mux)

	if err != nil {
		panic(err)
	}
}
