package handler

import (
	"fmt"
	"net/http"
	"strconv"

	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/cymiam/metrics-store/internal/service"
	"github.com/go-chi/chi/v5"
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
		err := handler.metricService.UpdateCounter(r.Context(), metricName, int64(metricValue))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		handler.logger.Info("Update metric",
			zap.String("MetricName", metricName),
			zap.String("MetricType", metricType),
			zap.Int("MetricValue", int(metricValue)))
		w.WriteHeader(http.StatusOK)
	case "gauge":
		err := handler.metricService.UpdateGauge(r.Context(), metricName, metricValue)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
		metric, err := handler.metricService.GetMetric(r.Context(), metricName, "counter")
		if err != nil {
			handler.logger.Error("error updating metric", zap.String("metric", metric.String()), zap.Error(err))

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
		metric, err := handler.metricService.GetMetric(r.Context(), metricName, "gauge")
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

	metrics, err := handler.metricService.GetAll(r.Context())

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

func (handler *MetricHandler) HandleUpdateJSON(w http.ResponseWriter, r *http.Request) {

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

		err := handler.metricService.UpdateCounter(r.Context(), metricName, *metric.Delta)

		if err != nil {
			handler.logger.Error("error updating metric", zap.String("metric", metric.String()), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

		err := handler.metricService.UpdateGauge(r.Context(), metricName, *metric.Value)

		if err != nil {
			handler.logger.Error("error updating metric", zap.String("metric", metric.String()), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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

func (handler *MetricHandler) HandleGetMetricJSON(w http.ResponseWriter, r *http.Request) {
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
		counter, err := handler.metricService.GetMetric(r.Context(), metricName, "counter")
		if err != nil {
			handler.logger.Error("error getting metric", zap.String("metric", metric.String()), zap.Error(err))

			w.WriteHeader(http.StatusNotFound)
			return
		}
		metric.Delta = counter.Delta
		easyjson.MarshalToHTTPResponseWriter(metric, w)
	case "gauge":
		gauge, err := handler.metricService.GetMetric(r.Context(), metricName, "gauge")
		if err != nil {
			handler.logger.Error("error getting metric", zap.String("metric", metric.String()), zap.Error(err))

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

func (handler *MetricHandler) ProcessBatchJSON(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-type", "application/json; charset=utf-8")

	var metrics models.Metrics

	err := easyjson.UnmarshalFromReader(r.Body, &metrics)

	if err != nil {
		handler.logger.Error("error unmarhsall json", zap.Error(err))
		http.Error(w, "Ошибка чтения тела", http.StatusBadRequest)
		return
	}

	err = handler.metricService.ProcessBatch(r.Context(), metrics)

	if err != nil {
		handler.logger.Error("Error processing batch", zap.Error(err))
		http.Error(w, "Ошибка обработки", http.StatusInternalServerError)
		return
	}

	handler.logger.Info("process batch success", zap.Int("metric count", len(metrics)))

	w.WriteHeader(http.StatusOK)
}
