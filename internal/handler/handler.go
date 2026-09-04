package handler

import (
	m "github.com/cymiam/metrics-store/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type MainRouterParams struct {
	Logger        *zap.Logger
	HealthHandler *HealthHandler
	MetricHandler *MetricHandler
}

func NewRouter(params MainRouterParams) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(m.GzipDecompressMidlleware)
	r.Use(m.RequestLoggerMiddleware(params.Logger))

	r.Route("/update", func(r chi.Router) {
		r.Post("/", params.MetricHandler.HandleUpdateJson)
		r.Post("/{metric_type}/{metric_name}/{metric_value}", params.MetricHandler.HandleUpdate)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", params.MetricHandler.HandleGetMetricJson)
		r.Get("/{metric_type}/{metric_name}", params.MetricHandler.HandleGetMetric)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/", params.MetricHandler.HandleGetMetrics)
	})

	r.Get("/ping", params.HealthHandler.HandlePing)
	return r
}
