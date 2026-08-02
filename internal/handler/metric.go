package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cymiam/metircs-store/internal/service"
)

type MetricHandler struct {
	metricService *service.MetricService
}

func NewMetricHandler() *MetricHandler {
	return &MetricHandler{
		metricService: service.NewMetricService(),
	}
}

func (handler *MetricHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	metricType := r.PathValue("metric_type")
	metricName := r.PathValue("metric_name")
	metricValue, err := strconv.ParseFloat(r.PathValue("metric_value"), 64)
	fmt.Printf("Got request, %s, %s, %f\n", metricType, metricName, metricValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	switch metricType {
	case "counter":
		handler.metricService.UpdateCounter(metricName)
		w.WriteHeader(http.StatusOK)
	case "gauge":
		handler.metricService.UpdateGauge(metricName, metricValue)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
	}
}
