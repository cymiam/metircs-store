package agent

import (
	"testing"

	config "github.com/cymiam/metrics-store/internal/config/agent"
	"github.com/stretchr/testify/assert"
)

func TestAgent_PollRuntimeMetrics(t *testing.T) {
	a := NewAgent(config.AgentConfig{
		Addr:           "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	})
	got := a.PollRuntimeMetrics()
	assert.Equal(t, 27, len(got))
}

func TestAgent_PollRuntimeMetricsPresent(t *testing.T) {

	a := NewAgent(config.AgentConfig{
		Addr:           "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	})
	tests := []struct {
		name  string
		names []string // description of this test case
	}{
		{
			name: "Metrics Present",
			names: []string{"Alloc",
				"BuckHashSys",
				"Frees",
				"GCCPUFraction",
				"GCSys",
				"HeapAlloc",
				"HeapIdle",
				"HeapInuse",
				"HeapObjects",
				"HeapReleased",
				"HeapSys",
				"LastGC",
				"Lookups",
				"MCacheInuse",
				"MCacheSys",
				"MSpanInuse",
				"MSpanSys",
				"Mallocs",
				"NextGC",
				"NumForcedGC",
				"NumGC",
				"OtherSys",
				"PauseTotalNs",
				"StackInuse",
				"StackSys",
				"Sys",
				"TotalAlloc"},
		},
	}

	metrics := a.PollRuntimeMetrics()
	got := make([]string, 0)

	for _, metric := range metrics {
		got = append(got, metric.ID)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, name := range tt.names {
				assert.Contains(t, got, name)
			}
		})
	}
}

func TestAgent_UpdatePollCount(t *testing.T) {

	a := NewAgent(config.AgentConfig{
		Addr:           "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	})
	a.PollRuntimeMetrics()
	assert.Equal(t, int64(1), a.PollCount)
}
