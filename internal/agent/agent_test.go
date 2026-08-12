package agent

import (
	"testing"

	config "github.com/cymiam/metircs-store/internal/config/agent"
	"github.com/stretchr/testify/assert"
)

func TestAgent_PollRuntimeMetrics(t *testing.T) {
	a := NewAgent(config.AgentConfig{
		Addr:           "localhost:8080",
		PollInterval:   2,
		ReportInterval: 10,
	})
	got := a.PollRuntimeMetrics()
	tests := []struct {
		name string // description of this test case
		want int64
	}{
		{
			name: "Test metric count",
			want: 27,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: update the condition below to compare got with tt.want.
			if len(got) != int(tt.want) {
				t.Errorf("PollRuntimeMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
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
	tests := []struct {
		name string // description of this test case
		want int64
	}{
		{
			name: "Test Update Pollcount",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.PollCount
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("PollRuntimeMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
