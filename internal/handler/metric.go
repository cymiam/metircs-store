package handler

import (
	"fmt"
	"net/http"
	"strconv"

	m "github.com/cymiam/metrics-store/internal/middleware"
	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type MetricHandler struct {
	metricService *service.MetricService
	logger        *zap.Logger
}

func NewMetricHandler(metricService *service.MetricService, logger *zap.Logger) *MetricHandler {
	return &MetricHandler{
		metricService: metricService,
		logger:        logger,
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
		handler.logger.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Int("MetricValue", int(metricValue)))
		w.WriteHeader(http.StatusOK)
	case "gauge":
		handler.metricService.UpdateGauge(metricName, metricValue)
		handler.logger.Info("Update metric",
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
		metric, err := handler.metricService.GetMetric(metricName, "counter")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value, err := metric.MetricValue()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
	case "gauge":
		metric, err := handler.metricService.GetMetric(metricName, "gauge")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		value, err := metric.MetricValue()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
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

	metrics, err := handler.metricService.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, metric := range metrics {

		value, err := metric.MetricValue()
		if err != nil {
			handler.logger.Error("html table", zap.Error(err))
		}
		body += "<tr>"
		body += fmt.Sprintf("<td>%s</td>", metric.ID)
		body += fmt.Sprintf("<td>%s</td>", value)
		body += "</tr>"
	}

	body += "</table>"
	w.Write([]byte(body))
}

func (handler *MetricHandler) HandleUpdateJson(w http.ResponseWriter, r *http.Request) {

	metric := models.Metric{}
	if err := easyjson.UnmarshalFromReader(r.Body, &metric); err != nil {
		handler.logger.Error("Error unmarhsalling json", zap.Error(err))
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

		if metric.Delta == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler.metricService.UpdateCounter(metricName, *metric.Delta)
		handler.logger.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Int("MetricValue", int(*metric.Delta)))
		w.WriteHeader(http.StatusOK)
	case "gauge":

		if metric.Value == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		handler.metricService.UpdateGauge(metricName, *metric.Value)
		handler.logger.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Float64("MetricValue", *metric.Value))
		w.WriteHeader(http.StatusOK)
	default:
		handler.logger.Info("Update metric failed",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType))
		http.Error(w, fmt.Sprintf("Неизвестный тип метрики: %s", metricType), http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HandleGetMetricJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	metric := models.Metric{}
	if err := easyjson.UnmarshalFromReader(r.Body, &metric); err != nil {
		handler.logger.Error("Error unmarhsalling json", zap.Error(err))
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
		counter, err := handler.metricService.GetMetric(metricName, "counter")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Delta = counter.Delta
		easyjson.MarshalToHTTPResponseWriter(metric, w)
	case "gauge":
		gauge, err := handler.metricService.GetMetric(metricName, "gauge")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Value = gauge.Value
		easyjson.MarshalToHTTPResponseWriter(metric, w)
	default:
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		http.Error(w, fmt.Sprintf("Неизвестный тип метрики: %s", metricType), http.StatusBadRequest)
	}
}

func (handler *MetricHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	err := handler.metricService.PingDB()
	if err != nil {
		handler.logger.Error("Error ping db", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewMetricRouter(metricHandler *MetricHandler, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(m.GzipDecompressMidlleware)
	r.Use(m.RequestLoggerMiddleware(logger))

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

	r.Get("/ping", metricHandler.HandleHealth)

	return r
}
