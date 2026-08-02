package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cymiam/metircs-store/internal/service"
	"github.com/go-chi/chi/v5"
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
	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")
	metricValue, err := strconv.ParseFloat(chi.URLParam(r, "metric_value"), 64)
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	switch metricType {
	case "counter":
		handler.metricService.UpdateCounter(metricName, int64(metricValue))
		fmt.Printf("Create metric: %s, %s, %f\n", metricType, metricName, metricValue)
		w.WriteHeader(http.StatusOK)
	case "gauge":
		handler.metricService.UpdateGauge(metricName, metricValue)
		fmt.Printf("Create metric: %s, %s, %f\n", metricType, metricName, metricValue)
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HandleGetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	switch metricType {
	case "counter":
		value, ok := handler.metricService.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		last := value[len(value)-1]
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%d", last)))
	case "gauge":
		value, ok := handler.metricService.GetGauge(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", value)))
	default:
		http.Error(w, "Неизвестный тип метрики", http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HanldeGetMetrics(w http.ResponseWriter, r *http.Request) {
	body := "======Metrics======\n"
	body += ("======Counters======\n")
	for k, v := range handler.metricService.GetCounters() {
		body += fmt.Sprintf("%s\t\t%d\n", k, v)
	}
	body += "======Gauges======\n"
	for k, v := range handler.metricService.GetGauges() {
		body += fmt.Sprintf("%s\t\t%v\n", k, v)
	}

	w.Write([]byte(body))
}

func MetricRouter() chi.Router {
	r := chi.NewRouter()
	metricHandler := NewMetricHandler()
	r.Route("/update", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Post("/{metric_type}/{metric_name}/{metric_value}", metricHandler.HandleUpdate)
		})
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metric_type}/{metric_name}", metricHandler.HandleGetMetric)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/", metricHandler.HanldeGetMetrics)
	})
	return r
}
