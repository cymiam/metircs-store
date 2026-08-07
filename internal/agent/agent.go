package agent

import (
	"net/http"
	"runtime"
	"time"

	config "github.com/cymiam/metircs-store/internal/config/agent"
	models "github.com/cymiam/metircs-store/internal/model"
)

type Agent struct {
	PollCount int64
	Client    http.Client
	Config    config.AgentConfig
}

func NewAgent() *Agent {
	return &Agent{
		PollCount: 0,
		Client: http.Client{
			Timeout: time.Second * 1, // интервал ожидания: 1 секунда
		},
		Config: config.ParseAgentConfig(),
	}
}

func (a *Agent) PollRuntimeMetrics() []models.Metric[float64] {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := []models.Metric[float64]{models.NewMetric("Alloc", float64(m.Alloc))}
	metrics = append(metrics, models.NewMetric("BuckHashSys", float64(m.BuckHashSys)))
	metrics = append(metrics, models.NewMetric("Frees", float64(m.Frees)))
	metrics = append(metrics, models.NewMetric("GCCPUFraction", m.GCCPUFraction))
	metrics = append(metrics, models.NewMetric("GCSys", float64(m.GCSys)))
	metrics = append(metrics, models.NewMetric("HeapAlloc", float64(m.HeapAlloc)))
	metrics = append(metrics, models.NewMetric("HeapIdle", float64(m.HeapIdle)))
	metrics = append(metrics, models.NewMetric("HeapInuse", float64(m.HeapInuse)))
	metrics = append(metrics, models.NewMetric("HeapObjects", float64(m.HeapObjects)))
	metrics = append(metrics, models.NewMetric("HeapReleased", float64(m.HeapReleased)))
	metrics = append(metrics, models.NewMetric("HeapSys", float64(m.HeapSys)))
	metrics = append(metrics, models.NewMetric("LastGC", float64(m.LastGC)))
	metrics = append(metrics, models.NewMetric("Lookups", float64(m.Lookups)))
	metrics = append(metrics, models.NewMetric("MCacheInuse", float64(m.MCacheInuse)))
	metrics = append(metrics, models.NewMetric("MCacheSys", float64(m.MCacheSys)))
	metrics = append(metrics, models.NewMetric("MSpanSys", float64(m.MSpanSys)))
	metrics = append(metrics, models.NewMetric("Mallocs", float64(m.Mallocs)))
	metrics = append(metrics, models.NewMetric("NextGC", float64(m.NextGC)))
	metrics = append(metrics, models.NewMetric("NumForcedGC", float64(m.NumForcedGC)))
	metrics = append(metrics, models.NewMetric("NumGC", float64(m.NumGC)))
	metrics = append(metrics, models.NewMetric("OtherSys", float64(m.OtherSys)))
	metrics = append(metrics, models.NewMetric("PauseTotalNs", float64(m.PauseTotalNs)))
	metrics = append(metrics, models.NewMetric("StackInuse", float64(m.StackInuse)))
	metrics = append(metrics, models.NewMetric("StackSys", float64(m.StackSys)))
	metrics = append(metrics, models.NewMetric("Sys", float64(m.Sys)))
	metrics = append(metrics, models.NewMetric("TotalAlloc", float64(m.TotalAlloc)))

	a.PollCount++
	return metrics
}
