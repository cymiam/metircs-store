package agent

import (
	"runtime"

	config "github.com/cymiam/metrics-store/internal/config/agent"
	models "github.com/cymiam/metrics-store/internal/model"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
	PollCount int64
	Client    resty.Client
	Config    config.AgentConfig
}

func NewAgent(cfg config.AgentConfig) *Agent {

	return &Agent{
		PollCount: 0,
		Client:    *resty.New(),
		Config:    cfg,
	}
}

func (a *Agent) PollRuntimeMetrics() []models.Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := []models.Metric{models.Metric{ID: "Alloc", Value: helper(float64(m.Alloc)), MType: "gauge"}}

	metrics = append(metrics, models.Metric{ID: "BuckHashSys", Value: helper(float64(m.BuckHashSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Frees", Value: helper(float64(m.Frees)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "GCCPUFraction", Value: helper(float64(m.GCCPUFraction)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "GCSys", Value: helper(float64(m.GCSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapAlloc", Value: helper(float64(m.HeapAlloc)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapIdle", Value: helper(float64(m.HeapIdle)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapInuse", Value: helper(float64(m.HeapInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapObjects", Value: helper(float64(m.HeapObjects)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapReleased", Value: helper(float64(m.HeapReleased)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "HeapSys", Value: helper(float64(m.HeapSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "LastGC", Value: helper(float64(m.LastGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Lookups", Value: helper(float64(m.Lookups)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MCacheInuse", Value: helper(float64(m.MCacheInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MCacheSys", Value: helper(float64(m.MCacheSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MSpanInuse", Value: helper(float64(m.MSpanInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "MSpanSys", Value: helper(float64(m.MSpanSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Mallocs", Value: helper(float64(m.Mallocs)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "NextGC", Value: helper(float64(m.NextGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "NumForcedGC", Value: helper(float64(m.NumForcedGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "NumGC", Value: helper(float64(m.NumGC)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "OtherSys", Value: helper(float64(m.OtherSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "PauseTotalNs", Value: helper(float64(m.PauseTotalNs)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "StackInuse", Value: helper(float64(m.StackInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "StackSys", Value: helper(float64(m.StackSys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "Sys", Value: helper(float64(m.Sys)), MType: "gauge"})
	metrics = append(metrics, models.Metric{ID: "TotalAlloc", Value: helper(float64(m.TotalAlloc)), MType: "gauge"})

	a.PollCount++
	return metrics
}

func helper(f float64) *float64 {
	return &f
}
