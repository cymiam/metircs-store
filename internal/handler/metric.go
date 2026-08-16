package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cymiam/metircs-store/internal/logger"
	models "github.com/cymiam/metircs-store/internal/model"
	"github.com/cymiam/metircs-store/internal/service"
	pkg "github.com/cymiam/metircs-store/pkg/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type MetricHandler struct {
	metricService *service.MetricService
}

func NewMetricHandler(metricService *service.MetricService) *MetricHandler {
	return &MetricHandler{
		metricService: metricService,
	}
}

func (handler *MetricHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")
	metricValue, err := strconv.ParseFloat(chi.URLParam(r, "metric_value"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	switch metricType {
	case "counter":
		handler.metricService.UpdateCounter(metricName, int64(metricValue))
		logger.Log.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Int("MetricValue", int(metricValue)))
		w.WriteHeader(http.StatusOK)
	case "gauge":
		handler.metricService.UpdateGauge(metricName, metricValue)
		logger.Log.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Float64("MetricValue", metricValue))
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, fmt.Sprintf("Неизвестный тип метрики: %s", metricType), http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HandleGetMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
		http.Error(w, fmt.Sprintf("Неизвестный тип метрики: %s", metricType), http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HandleGetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html; charset=utf-8")
	body := `
	<table>
	<tr>
		<th>Name</th>
		<th>Value</th>
	</tr>`
	for k, v := range handler.metricService.GetCounters() {
		body += "<tr>"
		body += fmt.Sprintf("<td>%s</td>", k)
		body += fmt.Sprintf("<td>%d</td>", v)
		body += "</tr>"
	}
	for k, v := range handler.metricService.GetGauges() {
		body += "<tr>"
		body += fmt.Sprintf("<td>%s</td>", k)
		body += fmt.Sprintf("<td>%f</td>", v)
		body += "</tr>"
	}
	body += "</table>"
	w.Write([]byte(body))
}

func (handler *MetricHandler) HandleUpdateJson(w http.ResponseWriter, r *http.Request) {

	metric := models.Metrics{}
	if err := easyjson.UnmarshalFromReader(r.Body, &metric); err != nil {
		logger.Log.Error("Error unmarhsalling json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricType := metric.MType
	metricName := metric.ID

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metricType {
	case "counter":
		handler.metricService.UpdateCounter(metricName, *metric.Delta)
		logger.Log.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Int("MetricValue", int(*metric.Delta)))
		w.WriteHeader(http.StatusOK)
	case "gauge":
		handler.metricService.UpdateGauge(metricName, *metric.Value)
		logger.Log.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Float64("MetricValue", *metric.Value))
		w.WriteHeader(http.StatusOK)
	default:
		logger.Log.Info("Update metric failed",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType))
		http.Error(w, fmt.Sprintf("Неизвестный тип метрики: %s", metricType), http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HandleGetMetricJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	metric := models.Metrics{}
	if err := easyjson.UnmarshalFromReader(r.Body, &metric); err != nil {
		logger.Log.Error("Error unmarhsalling json", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metricType := metric.MType
	metricName := metric.ID

	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metricType {
	case "counter":
		value, ok := handler.metricService.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		last := value[len(value)-1]
		metric.Delta = &last
		easyjson.MarshalToHTTPResponseWriter(metric, w)
	case "gauge":
		value, ok := handler.metricService.GetGauge(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Value = &value
		easyjson.MarshalToHTTPResponseWriter(metric, w)
	default:
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		http.Error(w, fmt.Sprintf("Неизвестный тип метрики: %s", metricType), http.StatusBadRequest)
	}
}

func NewMetricRouter(metricHandler *MetricHandler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(pkg.GzipDecompressMidlleware)
	r.Use(logger.RequestLogger)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", metricHandler.HandleUpdateJson)
		r.Post("/{metric_type}/{metric_name}/{metric_value}", metricHandler.HandleUpdate)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", metricHandler.HandleGetMetricJson)
		r.Get("/{metric_type}/{metric_name}", metricHandler.HandleGetMetric)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/", metricHandler.HandleGetMetrics)
	})

	return r
}
