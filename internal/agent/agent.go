package agent

import (
	"runtime"

	config "github.com/cymiam/metircs-store/internal/config/agent"
	models "github.com/cymiam/metircs-store/internal/model"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
	PollCount int64
	Client    resty.Client
	Config    config.AgentConfig
}

func NewAgent(configs ...config.AgentConfig) *Agent {
	var cfg config.AgentConfig

	if len(configs) > 0 {
		cfg = configs[0]
	} else {
		cfg = config.ParseAgentConfig()
	}

	return &Agent{
		PollCount: 0,
		Client:    *resty.New(),
		Config:    cfg,
	}
}

func (a *Agent) PollRuntimeMetrics() []models.Metrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := []models.Metrics{models.Metrics{ID: "Alloc", Value: helper(float64(m.Alloc)), MType: "gauge"}}

	metrics = append(metrics, models.Metrics{ID: "BuckHashSys", Value: helper(float64(m.BuckHashSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "Frees", Value: helper(float64(m.Frees)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "GCCPUFraction", Value: helper(float64(m.GCCPUFraction)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "GCSys", Value: helper(float64(m.GCSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "HeapAlloc", Value: helper(float64(m.HeapAlloc)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "HeapIdle", Value: helper(float64(m.HeapIdle)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "HeapInuse", Value: helper(float64(m.HeapInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "HeapObjects", Value: helper(float64(m.HeapObjects)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "HeapReleased", Value: helper(float64(m.HeapReleased)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "HeapSys", Value: helper(float64(m.HeapSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "LastGC", Value: helper(float64(m.LastGC)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "Lookups", Value: helper(float64(m.Lookups)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "MCacheInuse", Value: helper(float64(m.MCacheInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "MCacheSys", Value: helper(float64(m.MCacheSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "MSpanSys", Value: helper(float64(m.MSpanSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "Mallocs", Value: helper(float64(m.Mallocs)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "NextGC", Value: helper(float64(m.NextGC)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "NumForcedGC", Value: helper(float64(m.NumForcedGC)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "NumGC", Value: helper(float64(m.NumGC)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "OtherSys", Value: helper(float64(m.OtherSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "PauseTotalNs", Value: helper(float64(m.PauseTotalNs)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "StackInuse", Value: helper(float64(m.StackInuse)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "StackSys", Value: helper(float64(m.StackSys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "Sys", Value: helper(float64(m.Sys)), MType: "gauge"})
	metrics = append(metrics, models.Metrics{ID: "TotalAlloc", Value: helper(float64(m.TotalAlloc)), MType: "gauge"})

	a.PollCount++
	return metrics
}

func helper(f float64) *float64 {
	return &f
}
